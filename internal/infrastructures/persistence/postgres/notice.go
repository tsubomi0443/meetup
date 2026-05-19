package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetNoticeByQuestionIDs は質問 ID 一覧に紐づく通知を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - questionIDs []int64: 質問 ID 一覧
//
// return:
//   - []entity.Notice: 通知一覧
//   - error: DB エラー
func GetNoticeByQuestionIDs(ctx context.Context, db *gorm.DB, questionIDs []int64) (models []entity.Notice, err error) {
	if len(questionIDs) > 0 {
		models, err = gorm.G[entity.Notice](db).Where("question_id IN ?", questionIDs).Order("id").Find(ctx)
	}
	return
}

// GetNoticeByQuestion は質問に紐づく通知を1件取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - question entity.Question: 対象質問
//
// return:
//   - entity.Notice: 通知
//   - error: DB エラー
func GetNoticeByQuestion(ctx context.Context, db *gorm.DB, question entity.Question) (models entity.Notice, err error) {
	models, err = gorm.G[entity.Notice](db).Where("question_id = ?", question.ID).First(ctx)
	return
}

// GetNotice は通知一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.Notice: 通知一覧
//   - error: DB エラー
func GetNotice(ctx context.Context, db *gorm.DB) (models []entity.Notice, err error) {
	models, err = gorm.G[entity.Notice](db).
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

// GetNoticeByQuestionSilent は質問に紐づく通知をサイレントログで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - question entity.Question: 対象質問
//
// return:
//   - entity.Notice: 通知
//   - error: DB エラー
func GetNoticeByQuestionSilent(ctx context.Context, db *gorm.DB, question entity.Question) (model entity.Notice, err error) {
	model, err = gorm.G[entity.Notice](db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(logger.Silent),
	})).Where("question_id = ?", question.ID).First(ctx)
	return
}

// RegisterNoticeByQuestionID は回答期日接近の通知を質問 ID に紐づけて登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - questionID int64: 質問 ID
//
// return:
//   - error: DB エラー
func RegisterNoticeByQuestionID(ctx context.Context, db *gorm.DB, questionID int64) error {
	var content = "質問の回答期日が近づいています。"
	notice := entity.Notice{
		TypeID:     3,
		QuestionID: &questionID,
		Content:    &content,
	}
	return gorm.G[entity.Notice](db).Create(ctx, &notice)
}

// DeleteNoticeByID は指定 ID の通知を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 通知 ID
//
// return:
//   - error: DB エラー
func DeleteNoticeByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[entity.Notice](db).Where("id = ?", id).Delete(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteNoticeByQuestion は質問に紐づく通知を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - question entity.Question: 対象質問
//
// return:
//   - int64: 削除した通知 ID
//   - error: DB エラー
func DeleteNoticeByQuestion(ctx context.Context, db *gorm.DB, question entity.Question) (noticeID int64, err error) {
	n, err := GetNoticeByQuestionSilent(ctx, db, question)
	if err != nil {
		return -1, err
	}
	if err := DeleteNoticeByID(ctx, db, n.ID); err != nil {
		return -1, err
	}
	return n.ID, nil
}

// DeleteNoticeByQuestionID は質問 ID に紐づく通知を検索して削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - questionID int64: 質問 ID
//
// return:
//   - int64: 削除した通知 ID（見つからない場合は -1）
//   - error: DB エラー
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
