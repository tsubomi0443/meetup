package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetUserByID(ctx context.Context, db *gorm.DB, id int64) (model User, err error) {
	model, err = gorm.G[User](db).Where("id = ?", id).Select("id, name, email, memo, role_id").
		Preload("Role", commonPreloadBuilder()).
		First(ctx)
	return
}

func GetUserInfo(ctx context.Context, db *gorm.DB, emailOrName, pass string) (user User, err error) {
	user, err = gorm.G[User](db).Where("(email = ? OR name = ?) AND password = ?", emailOrName, emailOrName, pass).First(ctx)
	return
}

func GetUsers(ctx context.Context, db *gorm.DB) (models []User, err error) {
	models, err = gorm.G[User](db).
		Where("role_id <> ?", 1).
		Preload("Role", commonPreloadBuilder()).
		Not("role_id = 1").
		Select("id, name, email, memo, role_id").
		Order("id").
		Find(ctx)
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

// postgresTableCoalesceMaxID returns COALESCE(MAX(id), 0) for a whitelisted table name.
func postgresTableCoalesceMaxID(tx *gorm.DB, table string) (int64, error) {
	switch table {
	case "memos", "tag_managers", "related_questions":
	default:
		return 0, fmt.Errorf("postgresTableCoalesceMaxID: unsupported table %q", table)
	}
	var max int64
	q := fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", table)
	if err := tx.Raw(q).Scan(&max).Error; err != nil {
		return 0, err
	}
	return max, nil
}

// assignPostgreSQLBulkInsertZeros assigns sequential explicit primary keys to rows with ID==0 before
// GORM emits a multi-row INSERT. Otherwise one statement can mix RETURNING/Default nextval rows with rows
// that reuse client ids already present in MAX(id)...nextval bracket, producing memos_pkey duplicate key.
func assignMemoBulkInsertZeros(tx *gorm.DB, rows []Memo) error {
	hasZero := false
	batchMax := int64(0)
	for i := range rows {
		id := rows[i].ID
		if id == 0 {
			hasZero = true
		} else if id > batchMax {
			batchMax = id
		}
	}
	if !hasZero {
		return nil
	}
	dbMax, err := postgresTableCoalesceMaxID(tx, "memos")
	if err != nil {
		return err
	}
	next := batchMax
	if dbMax > next {
		next = dbMax
	}
	for i := range rows {
		if rows[i].ID == 0 {
			next++
			rows[i].ID = next
		}
	}
	return nil
}

func assignTagManagerBulkInsertZeros(tx *gorm.DB, rows []TagManager) error {
	hasZero := false
	batchMax := int64(0)
	for i := range rows {
		id := rows[i].ID
		if id == 0 {
			hasZero = true
		} else if id > batchMax {
			batchMax = id
		}
	}
	if !hasZero {
		return nil
	}
	dbMax, err := postgresTableCoalesceMaxID(tx, "tag_managers")
	if err != nil {
		return err
	}
	next := batchMax
	if dbMax > next {
		next = dbMax
	}
	for i := range rows {
		if rows[i].ID == 0 {
			next++
			rows[i].ID = next
		}
	}
	return nil
}

func assignRelatedQuestionBulkInsertZeros(tx *gorm.DB, rows []RelatedQuestion) error {
	hasZero := false
	batchMax := int64(0)
	for i := range rows {
		id := rows[i].ID
		if id == 0 {
			hasZero = true
		} else if id > batchMax {
			batchMax = id
		}
	}
	if !hasZero {
		return nil
	}
	dbMax, err := postgresTableCoalesceMaxID(tx, "related_questions")
	if err != nil {
		return err
	}
	next := batchMax
	if dbMax > next {
		next = dbMax
	}
	for i := range rows {
		if rows[i].ID == 0 {
			next++
			rows[i].ID = next
		}
	}
	return nil
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
		// 4) タグ紐づけ（tag_managers）— has_many Replace は子の FK を NULL 更新するため
		// NOT NULL 制約（question_id）と相性が悪い。DELETE + INSERT で置き換える。
		var tagRows []TagManager
		for _, tm := range q.TagManagers {
			if tm.TagID == 0 {
				continue
			}
			tagRows = append(tagRows, TagManager{
				ID:         tm.ID,
				TagID:      tm.TagID,
				QuestionID: q.ID,
			})
		}
		if err := tx.Unscoped().Where("question_id = ?", q.ID).Delete(&TagManager{}).Error; err != nil {
			return err
		}
		if len(tagRows) > 0 {
			if err := assignTagManagerBulkInsertZeros(tx, tagRows); err != nil {
				return err
			}
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
			memo := Memo{
				ID:         m.ID,
				UserID:     m.UserID,
				Content:    content,
				QuestionID: q.ID,
			}
			memoRows = append(memoRows, memo)
		}
		if err := tx.Unscoped().Where("question_id = ?", q.ID).Delete(&Memo{}).Error; err != nil {
			return err
		}
		if len(memoRows) > 0 {
			if err := assignMemoBulkInsertZeros(tx, memoRows); err != nil {
				return err
			}
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
				ID:                rq.ID,
				QuestionID:        q.ID,
				RelatedQuestionID: rq.RelatedQuestionID,
			})
		}
		if err := tx.Unscoped().Where("question_id = ?", q.ID).Delete(&RelatedQuestion{}).Error; err != nil {
			return err
		}
		if len(relatedRows) > 0 {
			if err := assignRelatedQuestionBulkInsertZeros(tx, relatedRows); err != nil {
				return err
			}
			if err := tx.Create(&relatedRows).Error; err != nil {
				return err
			}
		}
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
		Preload("Answer", commonPreloadBuilder()).
		Preload("Answer.User", commonPreloadBuilder()).
		Preload("Answer.User.Role", commonPreloadBuilder()).
		Preload("Answer.ReferManagers", commonPreloadBuilder()).
		Preload("Answer.ReferManagers.Refer", commonPreloadBuilder()).
		Preload("Memos", commonPreloadBuilder()).
		Preload("Memos.User", commonPreloadBuilder()).
		Preload("Memos.User.Role", commonPreloadBuilder()).
		Preload("TagManagers", commonPreloadBuilder()).
		Preload("TagManagers.Tag", commonPreloadBuilder()).
		Preload("TagManagers.Tag.Category", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", commonPreloadBuilder()).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

func DeleteUserByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[User](db).
		Preload("Role", commonPreloadBuilder()).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}

func GetQuestion(ctx context.Context, db *gorm.DB, id int64) (model Question, err error) {
	model, err = gorm.G[Question](db).
		Preload("Answer", commonPreloadBuilder()).
		Preload("Answer.User", userPreloadBuilder(false)).
		Preload("Answer.User.Role", commonPreloadBuilder()).
		Preload("Answer.ReferManagers", commonPreloadBuilder()).
		Preload("Answer.ReferManagers.Refer", commonPreloadBuilder()).
		Preload("Memos", commonPreloadBuilder()).
		Preload("Memos.User", userPreloadBuilder(false)).
		Preload("Memos.User.Role", commonPreloadBuilder()).
		Preload("TagManagers", commonPreloadBuilder()).
		Preload("TagManagers.Tag", commonPreloadBuilder()).
		Preload("TagManagers.Tag.Category", commonPreloadBuilder()).
		Preload("RelatedQuestions", commonPreloadBuilder()).
		Preload("RelatedQuestions.RelatedQuestion", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", userPreloadBuilder(false)).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Preload("RelatedQuestions", commonPreloadBuilder()).
		Preload("RelatedQuestions.RelatedQuestion", commonPreloadBuilder()).
		Where("id = ?", id).
		First(ctx)
	return
}

func GetQuestions(ctx context.Context, db *gorm.DB) (models []Question, err error) {
	models, err = gorm.G[Question](db).
		Preload("Answer", commonPreloadBuilder()).
		Preload("Answer.User", userPreloadBuilder(false)).
		Preload("Answer.User.Role", commonPreloadBuilder()).
		Preload("Answer.ReferManagers", commonPreloadBuilder()).
		Preload("Answer.ReferManagers.Refer", commonPreloadBuilder()).
		Preload("Memos", commonPreloadBuilder()).
		Preload("Memos.User", userPreloadBuilder(false)).
		Preload("Memos.User.Role", commonPreloadBuilder()).
		Preload("TagManagers", commonPreloadBuilder()).
		Preload("TagManagers.Tag", commonPreloadBuilder()).
		Preload("TagManagers.Tag.Category", commonPreloadBuilder()).
		Preload("RelatedQuestions", commonPreloadBuilder()).
		Preload("RelatedQuestions.RelatedQuestion", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", userPreloadBuilder(false)).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Preload("RelatedQuestions", commonPreloadBuilder()).
		Preload("RelatedQuestions.RelatedQuestion", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

func GetTags(ctx context.Context, db *gorm.DB) (models []Tag, err error) {
	models, err = gorm.G[Tag](db).
		Preload("Category", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

func GetTagByID(ctx context.Context, db *gorm.DB, id int64) (models Tag, err error) {
	models, err = gorm.G[Tag](db).
		Preload("Category", commonPreloadBuilder()).
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
		models, err = gorm.G[Notice](db).Where("question_id IN ?", questionIDs).Order("id").Find(ctx)
	}
	return
}

func GetNoticeByQuestion(ctx context.Context, db *gorm.DB, question Question) (models Notice, err error) {
	models, err = gorm.G[Notice](db).Where("question_id = ?", question.ID).First(ctx)
	return
}

func GetNotice(ctx context.Context, db *gorm.DB) (models []Notice, err error) {
	models, err = gorm.G[Notice](db).
		Preload("NoticeType", commonPreloadBuilder()).
		Preload("Question", commonPreloadBuilder()).
		Preload("Question.Support", commonPreloadBuilder()).
		Preload("Question.Support.SupportStatus", commonPreloadBuilder()).
		Preload("Question.TagManagers", commonPreloadBuilder()).
		Preload("Question.TagManagers.Tag", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

func GetNoticeByQuestionSilent(ctx context.Context, db *gorm.DB, question Question) (model Notice, err error) {
	model, err = gorm.G[Notice](db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(logger.Silent),
	})).Where("question_id = ?", question.ID).First(ctx)
	return
}

func RegisterNoticeByQuestionID(ctx context.Context, db *gorm.DB, questionID int64) error {
	var content = "質問の回答期日が近づいています。"
	notice := Notice{
		TypeID:       3,
		QuestionID:   &questionID,
		Content:      &content,
	}
	return gorm.G[Notice](db).Create(ctx, &notice)
}

func DeleteNoticeByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[Notice](db).Where("id = ?", id).Delete(ctx); err != nil {
		return err
	}
	return nil
}

func DeleteNoticeByQuestion(ctx context.Context, db *gorm.DB, question Question) (noticeID int64, err error) {
	n, err := GetNoticeByQuestionSilent(ctx, db, question)
	if err != nil {
		return -1, err
	}
	if err := DeleteNoticeByID(ctx, db, n.ID); err != nil {
		return -1, err
	}
	return n.ID, nil
}

func DeleteNoticeByQuestionID(ctx context.Context, db *gorm.DB, questionID int64) (deletedID int64, err error) {
	notices, err := GetNotice(ctx, db)
	if err != nil {
		return -1, err
	}
	for _, n := range notices {
		if n.QuestionID != nil && *n.QuestionID == questionID {
			if err := DeleteNoticeByID(ctx, db, n.ID); err != nil {
				return n.ID, err
			}
			return n.ID, nil
		}
	}
	return -1, gorm.ErrRecordNotFound
}

func GetMaxByColumn[T any](ctx context.Context, db *gorm.DB, columnName string) int64 {
	var max sql.NullInt64
	err := db.WithContext(ctx).Model(new(T)).
		Select(fmt.Sprintf("MAX(%s)", columnName)).
		Take(&max).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return -1
	}

	if !max.Valid {
		return -1
	}
	return max.Int64

}

func commonPreloadBuilder() func(db gorm.PreloadBuilder) error {
	return func(db gorm.PreloadBuilder) error {
		db.Order("id")
		return nil
	}
}

func userPreloadBuilder(includePassword bool) func(db gorm.PreloadBuilder) error {
	return func(db gorm.PreloadBuilder) error {
		if includePassword {
			db.Select("id, name, email, password, memo, role_id")
		} else {
			db.Select("id, name, email, memo, role_id")
		}
		db.Order("id ASC")
		return nil
	}
}
