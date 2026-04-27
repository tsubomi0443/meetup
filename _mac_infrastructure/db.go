package infrastructure

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetUserByID(ctx context.Context, db *gorm.DB, id int64) (model User, err error) {
	model, err = gorm.G[User](db).Where("id = ?", id).Select("id, name, email, memo, role_id").
		Preload("Role", nil).
		First(ctx)
	return
}

func GetUserInfo(ctx context.Context, db *gorm.DB, email, pass string) (user User, err error) {
	user, err = gorm.G[User](db).Where("email = ? AND password = ?", email, pass).First(ctx)
	return
}

func GetUsers(ctx context.Context, db *gorm.DB) (models []User, err error) {
	models, err = gorm.G[User](db).
		Where("role_id <> ?", 1).
		Preload("Role", nil).
		Not("role_id = 1").
		Select("id, name, email, memo, role_id").Find(ctx)
	return
}

func Register[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) error {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Create(ctx, &model)
}

func Updates[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) (int, error) {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Updates(ctx, model)
}

// UpdateByID は主キー id を WHERE に固定して単一モデルを更新する。
// 他の更新系メソッドに合わせて gorm.G[T](db) の Updates を利用し、トランザクションは張らない。
// omit は generic ビルダの Omit にそのまま渡す（関連名・列名どちらも可）。
func UpdateByID[T any](ctx context.Context, db *gorm.DB, id int64, model T, omit ...string) (int, error) {
	return gorm.G[T](db.WithContext(ctx)).
		Omit(omit...).
		Where("id = ?", id).
		Updates(ctx, model)
}

// updateInTransaction は単一の DB トランザクション内で gorm.Updates を実行する。
// omit に関連名を渡し、Role/Category など中間テーブル向けの関連を更新対象から外す。
// 既存の updates と違い、こちらは更新用。事前読み込みは行わない。
func UpdateInTransaction[T any](ctx context.Context, db *gorm.DB, model T, omit ...string) (rowsAffected int, err error) {
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := model
		res := tx.Omit(omit...).Model(&m).Updates(&m)
		rowsAffected = int(res.RowsAffected)
		return res.Error
	})
	return
}

// updateQuestionInTransaction は Question の1行更新と、QuestionToEntity で組み立てた
// 1対多の関連（Answer, Support, TagManagers, Memos, RelatedQuestions）の同期を1トランザクションで行う。
// フォーム上の下位の質問（SubQuestions）やエスカレーション等は本関数では永続化しない。
func UpdateQuestionInTransaction(ctx context.Context, db *gorm.DB, q Question) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Answer: 新規作成 or 既存更新し、親の answer_id を整合させる
		if q.Answer != nil {
			a := *q.Answer
			if a.ID == 0 {
				if err := tx.Create(&a).Error; err != nil {
					return err
				}
				aid := a.ID
				q.AnswerID = &aid
				q.Answer = &a
			} else {
				if _, err := gorm.G[Answer](tx).
					Omit("User", "ReferManagers").
					Where("id = ?", a.ID).
					Updates(ctx, a); err != nil {
					return err
				}
				if q.AnswerID == nil {
					aid := a.ID
					q.AnswerID = &aid
				}
			}
		}
		// 2) Support: 新規作成 or 既存更新
		if q.Support != nil {
			s := *q.Support
			if s.ID == 0 {
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
				sid := s.ID
				q.SupportID = &sid
				q.Support = &s
			} else {
				if err := tx.Model(&s).Omit("User", "SupportStatus").Updates(&s).Error; err != nil {
					return err
				}
				if q.SupportID == nil {
					sid := s.ID
					q.SupportID = &sid
					q.Support = &s
				}
			}
		}
		// 3) 親の questions: スカラ列と FK のみ（CreatedAt は更新で変えない）
		if _, err := gorm.G[Question](tx).
			Omit("Answer", "Memos", "Notices", "TagManagers", "Support", "RelatedQuestions", "CreatedAt").
			Where("id = ?", q.ID).
			Updates(ctx, q); err != nil {
			return err
		}
		// ref := Question{ID: q.ID}
		// 4) タグ紐づけ（tag_managers）— has_many Replace は子の FK を NULL 更新するため
		// NOT NULL 制約（question_id）と相性が悪い。DELETE + INSERT で置き換える。
		var tagRows []TagManager
		for _, tm := range q.TagManagers {
			if tm.TagID == 0 {
				continue
			}
			tagRows = append(tagRows, TagManager{
				ID:         0,
				TagID:      tm.TagID,
				QuestionID: q.ID,
			})
		}
		if err := tx.Where("question_id = ?", q.ID).Delete(&TagManager{}).Error; err != nil {
			return err
		}
		if len(tagRows) > 0 {
			if err := tx.Create(&tagRows).Error; err != nil {
				return err
			}
		}
		// 5) メモ（memos）— 同じ理由で Association.Replace ではなく DELETE + INSERT
		var memoRows []Memo
		for _, m := range q.Memos {
			content := strings.TrimSpace(m.Content)
			if m.UserID == 0 || content == "" {
				continue
			}
			memoRows = append(memoRows, Memo{
				ID:         0,
				UserID:     m.UserID,
				Content:    content,
				QuestionID: q.ID,
			})
		}
		if err := tx.Where("question_id = ?", q.ID).Delete(&Memo{}).Error; err != nil {
			return err
		}
		if len(memoRows) > 0 {
			if err := tx.Create(&memoRows).Error; err != nil {
				return err
			}
		}
		// 6) 関連質問（related_questions）— tag_managers / memos と同様に DELETE + INSERT
		var relatedRows []RelatedQuestion
		for _, rq := range q.RelatedQuestions {
			if rq.RelatedQuestionID == 0 || rq.RelatedQuestionID == q.ID {
				continue
			}
			relatedRows = append(relatedRows, RelatedQuestion{
				ID:                0,
				QuestionID:        q.ID,
				RelatedQuestionID: rq.RelatedQuestionID,
			})
		}
		if err := tx.Where("question_id = ?", q.ID).Delete(&RelatedQuestion{}).Error; err != nil {
			return err
		}
		if len(relatedRows) > 0 {
			if err := tx.Create(&relatedRows).Error; err != nil {
				return err
			}
		}
		// // 7) 通知
		// for i := range q.Notices {
		// 	if q.Notices[i].QuestionID == nil {
		// 		qid := q.ID
		// 		q.Notices[i].QuestionID = &qid
		// 	}
		// }
		// if len(q.Notices) == 0 {
		// 	if err := tx.Model(&ref).Association("Notices").Clear(); err != nil {
		// 		return err
		// 	}
		// } else {
		// 	if err := tx.Model(&ref).Association("Notices").Replace(q.Notices); err != nil {
		// 		return err
		// 	}
		// }
		return nil
	})
}

func DeleteQuestionByID(ctx context.Context, db *gorm.DB, id int64) error {
	if err := db.WithContext(ctx).
		Where("question_id = ? OR related_question_id = ?", id, id).
		Delete(&RelatedQuestion{}).Error; err != nil {
		return err
	}
	if _, err := gorm.G[Question](db).
		Preload("Answer", nil).
		Preload("Answer.User", nil).
		Preload("Answer.User.Role", nil).
		Preload("Answer.ReferManagers", nil).
		Preload("Answer.ReferManagers.Refer", nil).
		Preload("Memos", nil).
		Preload("Memos.User", nil).
		Preload("Memos.User.Role", nil).
		Preload("TagManagers", nil).
		Preload("TagManagers.Tag", nil).
		Preload("TagManagers.Tag.Category", nil).
		Preload("Support", nil).
		Preload("Support.User", nil).
		Preload("Support.User.Role", nil).
		Preload("Support.SupportStatus", nil).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

func DeleteUserByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[User](db).
		Preload("Role", nil).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

func GetQuestion(ctx context.Context, db *gorm.DB, id int64) (model Question, err error) {
	model, err = gorm.G[Question](db).
		Preload("Answer", nil).
		Preload("Answer.User", func(pb gorm.PreloadBuilder) error {
			pb.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Answer.User.Role", nil).
		Preload("Answer.ReferManagers", nil).
		Preload("Answer.ReferManagers.Refer", nil).
		Preload("Memos", nil).
		Preload("Memos.User", func(pb gorm.PreloadBuilder) error {
			pb.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Memos.User.Role", nil).
		Preload("TagManagers", nil).
		Preload("TagManagers.Tag", nil).
		Preload("TagManagers.Tag.Category", nil).
		Preload("RelatedQuestions", nil).
		Preload("RelatedQuestions.RelatedQuestion", nil).
		Preload("Support", nil).
		Preload("Support.User", func(pb gorm.PreloadBuilder) error {
			pb.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Support.User.Role", nil).
		Preload("Support.SupportStatus", nil).
		Where("id = ?", id).
		First(ctx)
	return
}

func GetQuestions(ctx context.Context, db *gorm.DB) (models []Question, err error) {
	models, err = gorm.G[Question](db).
		Preload("Answer", nil).
		Preload("Answer.User", func(pb gorm.PreloadBuilder) error {
			pb.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Answer.User.Role", nil).
		Preload("Answer.ReferManagers", nil).
		Preload("Answer.ReferManagers.Refer", nil).
		Preload("Memos", nil).
		Preload("Memos.User", func(pb gorm.PreloadBuilder) error {
			pb.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Memos.User.Role", nil).
		Preload("TagManagers", nil).
		Preload("TagManagers.Tag", nil).
		Preload("TagManagers.Tag.Category", nil).
		Preload("RelatedQuestions", nil).
		Preload("RelatedQuestions.RelatedQuestion", nil).
		Preload("Support", nil).
		Preload("Support.User", func(pb gorm.PreloadBuilder) error {
			pb.Select("id", "name", "email", "role_id")
			return nil
		}).
		Preload("Support.User.Role", nil).
		Preload("Support.SupportStatus", nil).
		Find(ctx)
	return
}

func GetTags(ctx context.Context, db *gorm.DB) (models []Tag, err error) {
	models, err = gorm.G[Tag](db).
		Preload("Category", nil).
		Find(ctx)
	return
}

func GetTagByID(ctx context.Context, db *gorm.DB, id int64) (models Tag, err error) {
	models, err = gorm.G[Tag](db).
		Preload("Category", nil).
		Where("id = ?", id).
		First(ctx)
	return
}

func DeleteTagByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[Tag](db).Where("id = ?", id).Limit(1).Delete(ctx); err != nil {
		return err
	}
	return nil
}

func GetNoticeByQuestionIDs(ctx context.Context, db *gorm.DB, questionIDs []int64) (models []Notice, err error) {
	if len(questionIDs) > 0 {
		models, err = gorm.G[Notice](db).Where("question_id IN ?", questionIDs).Find(ctx)
	}
	return
}

func GetNoticeByQuestion(ctx context.Context, db *gorm.DB, question Question) (models Notice, err error) {
	models, err = gorm.G[Notice](db).Where("question_id = ?", question.ID).First(ctx)
	return
}

func GetNotice(ctx context.Context, db *gorm.DB) (models []Notice, err error) {
	models, err = gorm.G[Notice](db).
		Preload("NoticeType", nil).
		Preload("Question", nil).
		Preload("Question.Support", nil).
		Preload("Question.Support.SupportStatus", nil).
		Preload("Question.TagManagers", nil).
		Preload("Question.TagManagers.Tag", nil).
		Find(ctx)
	return
}

func GetNoticeByQuestionSilent(ctx context.Context, db *gorm.DB, question Question) (models Notice, err error) {
	models, err = gorm.G[Notice](db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(logger.Silent),
	})).Where("question_id = ?", question.ID).First(ctx)
	return
}

func RegisterNoticeByQuestionID(ctx context.Context, db *gorm.DB, questionID int64) error {
	var content = "質問の回答期日が近づいています。"
	notice := Notice{
		TypeID:     3,
		QuestionID: &questionID,
		Content:    &content,
	}
	return gorm.G[Notice](db).Create(ctx, &notice)
}

func DeleteNoticeByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[Notice](db).Where("id = ?", id).Delete(ctx); err != nil {
		return err
	}
	return nil
}
