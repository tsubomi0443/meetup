package postgres

import (
	"context"
	"strings"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// UpdateQuestionInTransaction は質問本体と1対多関連を差分同期で1トランザクション更新する。
// related_questions のみ物理削除（ユニーク制約のため）。SubQuestions 等は永続化しない。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - q entity.Question: 更新後の質問エンティティ
//
// return:
//   - error: DB エラー
func UpdateQuestionInTransaction(ctx context.Context, db *gorm.DB, q entity.Question) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) entity.Support: 新規作成 or 既存更新、またはフォームに無ければ既存の support を detach
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
				if err := tx.Model(&s).Omit("User", "entity.SupportStatus").Updates(&s).Error; err != nil {
					return err
				}
				if q.SupportID == nil {
					sid := s.ID
					q.SupportID = &sid
					q.Support = &s
				}
			}
		} else {
			if err := DetachQuestionSupportTx(tx, q.ID); err != nil {
				return err
			}
			q.SupportID = nil
		}
		// 2) 親の questions: スカラ列と FK のみ（CreatedAt は更新で変えない）
		if _, err := gorm.G[entity.Question](tx).
			Omit("Answer", "Memos", "Notices", "TagManagers", "entity.Support", "RelatedQuestions", "SenderTalks", "TalkroomID", "CreatedAt").
			Where("id = ?", q.ID).
			Updates(ctx, q); err != nil {
			return err
		}
		// 3) タグ紐づけ（tag_managers）— 自然キー tag_id で差分同期、削除は論理削除
		var tagRows []entity.TagManager
		for _, tm := range q.TagManagers {
			if tm.TagID == 0 {
				continue
			}
			tagRows = append(tagRows, entity.TagManager{
				ID:         tm.ID,
				TagID:      tm.TagID,
				QuestionID: q.ID,
			})
		}
		if _, err := syncChildrenByKey(tx, "tag_managers", "question_id", q.ID, tagRows,
			func(t *entity.TagManager) int64 { return t.TagID },
			func(t *entity.TagManager) *int64 { return &t.ID },
			nil,
			true,
		); err != nil {
			return err
		}
		// 4) 回答（answers）— ID で差分同期、削除は論理削除。refer_managers は各回答で refer_id 自然キーで差分・論理削除。
		var answerRows []entity.Answer
		var referRowsPerAnswer [][]entity.ReferManager
		for _, a := range q.Answer {
			content := strings.TrimSpace(a.Content)
			if a.UserID == 0 || content == "" {
				continue
			}
			answerRows = append(answerRows, entity.Answer{
				ID:         a.ID,
				UserID:     a.UserID,
				Content:    content,
				IsFinal:    a.IsFinal,
				QuestionID: q.ID,
			})
			var refs []entity.ReferManager
			for _, rm := range a.ReferManagers {
				referID, err := resolveReferIDForSync(tx, rm)
				if err != nil {
					return err
				}
				if referID == 0 {
					continue
				}
				refs = append(refs, entity.ReferManager{ReferID: referID})
			}
			referRowsPerAnswer = append(referRowsPerAnswer, refs)
		}
		deletedAnswerIDs, err := syncChildrenByKey(tx, "answers", "question_id", q.ID, answerRows,
			func(a *entity.Answer) int64 { return a.ID },
			func(a *entity.Answer) *int64 { return &a.ID },
			func(tx *gorm.DB, prev entity.Answer, next *entity.Answer) error {
				return tx.Model(&entity.Answer{}).Where("id = ?", prev.ID).
					Updates(map[string]any{
						"content":  next.Content,
						"is_final": next.IsFinal,
						"user_id":  next.UserID,
					}).Error
			},
			true,
		)
		if err != nil {
			return err
		}
		if len(deletedAnswerIDs) > 0 {
			if err := tx.Where("answer_id IN ?", deletedAnswerIDs).Delete(&entity.ReferManager{}).Error; err != nil {
				return err
			}
		}
		for i := range answerRows {
			refs := referRowsPerAnswer[i]
			for j := range refs {
				refs[j].AnswerID = answerRows[i].ID
			}
			if _, err := syncChildrenByKey(tx, "refer_managers", "answer_id", answerRows[i].ID, refs,
				func(r *entity.ReferManager) int64 { return r.ReferID },
				func(r *entity.ReferManager) *int64 { return &r.ID },
				nil,
				true,
			); err != nil {
				return err
			}
		}
		// 5) メモ（memos）— ID で差分同期、削除は論理削除
		var memoRows []entity.Memo
		for _, m := range q.Memos {
			content := strings.TrimSpace(m.Content)
			if m.UserID == 0 || content == "" {
				continue
			}
			memoRows = append(memoRows, entity.Memo{
				ID:         m.ID,
				UserID:     m.UserID,
				Content:    content,
				QuestionID: q.ID,
			})
		}
		if _, err := syncChildrenByKey(tx, "memos", "question_id", q.ID, memoRows,
			func(m *entity.Memo) int64 { return m.ID },
			func(m *entity.Memo) *int64 { return &m.ID },
			func(tx *gorm.DB, prev entity.Memo, next *entity.Memo) error {
				return tx.Model(&entity.Memo{}).Where("id = ?", prev.ID).
					Updates(map[string]any{"content": next.Content, "user_id": next.UserID}).Error
			},
			true,
		); err != nil {
			return err
		}
		// 6) 関連質問（related_questions）— 自然キー related_question_id で差分同期。
		// doc/db/INIT.sql の uq_related_questions UNIQUE(question_id, related_question_id) により、
		// 論理削除（deleted_at をセット）してもユニーク制約が deleted_at を区別しないため、
		// 同じ自然キーで再追加すると衝突する。そのためここだけ物理削除（Unscoped）で同期する。
		var relatedRows []entity.RelatedQuestion
		for _, rq := range q.RelatedQuestions {
			if rq.RelatedQuestionID == 0 || rq.RelatedQuestionID == q.ID {
				continue
			}
			relatedRows = append(relatedRows, entity.RelatedQuestion{
				ID:                rq.ID,
				QuestionID:        q.ID,
				RelatedQuestionID: rq.RelatedQuestionID,
			})
		}
		if _, err := syncChildrenByKey(tx, "related_questions", "question_id", q.ID, relatedRows,
			func(r *entity.RelatedQuestion) int64 { return r.RelatedQuestionID },
			func(r *entity.RelatedQuestion) *int64 { return &r.ID },
			nil,
			false,
		); err != nil {
			return err
		}
		return nil
	})
}

// DetachQuestionSupportTx は質問の support_id を NULL にし、紐づく supports 行を削除する（1:1 前提）。
//
// args:
//   - tx *gorm.DB: トランザクション
//   - questionID int64: 質問 ID
//
// return:
//   - error: DB エラー
func DetachQuestionSupportTx(tx *gorm.DB, questionID int64) error {
	var current entity.Question
	if err := tx.Select("id", "support_id").
		Where("id = ?", questionID).
		Take(&current).Error; err != nil {
		return err
	}
	if err := tx.Model(&entity.Question{}).
		Where("id = ?", questionID).
		Update("support_id", nil).Error; err != nil {
		return err
	}
	if current.SupportID != nil && *current.SupportID != 0 {
		if err := tx.Unscoped().
			Where("id = ?", *current.SupportID).
			Delete(&entity.Support{}).Error; err != nil {
			return err
		}
	}
	return nil
}

// DeleteQuestionByID は指定 ID の質問と関連行を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 質問 ID
//
// return:
//   - error: DB エラー
func DeleteQuestionByID(ctx context.Context, db *gorm.DB, id int64) error {
	if err := db.WithContext(ctx).
		Where("question_id = ? OR related_question_id = ?", id, id).
		Delete(&entity.RelatedQuestion{}).Error; err != nil {
		return err
	}
	if _, err := gorm.G[entity.Question](db).
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

// GetQuestion は指定 ID の質問を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 質問 ID
//
// return:
//   - entity.Question: 質問
//   - error: DB エラー
func GetQuestion(ctx context.Context, db *gorm.DB, id int64) (model entity.Question, err error) {
	model, err = gorm.G[entity.Question](db).
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
		Preload("SenderTalks", commonPreloadBuilder()).
		Preload("SenderTalks.Sender", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", userPreloadBuilder(false)).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Where("id = ?", id).
		First(ctx)
	return
}

// GetQuestions は質問一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.Question: 質問一覧
//   - error: DB エラー
func GetQuestions(ctx context.Context, db *gorm.DB) (models []entity.Question, err error) {
	models, err = gorm.G[entity.Question](db).
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
		Preload("SenderTalks", commonPreloadBuilder()).
		Preload("SenderTalks.Sender", commonPreloadBuilder()).
		Preload("Support", commonPreloadBuilder()).
		Preload("Support.User", userPreloadBuilder(false)).
		Preload("Support.User.Role", commonPreloadBuilder()).
		Preload("Support.SupportStatus", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}
