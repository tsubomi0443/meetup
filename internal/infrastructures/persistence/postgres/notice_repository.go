package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// NoticeRepository は通知の永続化を担う。
type NoticeRepository struct {
	DB *gorm.DB
}

// NewNoticeRepository は NoticeRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *NoticeRepository: リポジトリ
func NewNoticeRepository(db *gorm.DB) *NoticeRepository {
	return &NoticeRepository{DB: db}
}

// GetAll は通知一覧を関連込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Notice: 通知一覧
//   - error: DB エラー
func (r *NoticeRepository) GetAll(ctx context.Context) ([]entity.Notice, error) {
	return GetNotice(ctx, r.DB)
}

// GetByQuestionSilent は質問に紐づく通知をサイレントログで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - question entity.Question: 対象質問
//
// return:
//   - entity.Notice: 通知
//   - error: DB エラー
func (r *NoticeRepository) GetByQuestionSilent(ctx context.Context, question entity.Question) (entity.Notice, error) {
	return GetNoticeByQuestionSilent(ctx, r.DB, question)
}

// RegisterByQuestionID は質問 ID に紐づく期限通知を登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - questionID int64: 質問 ID
//
// return:
//   - error: DB エラー
func (r *NoticeRepository) RegisterByQuestionID(ctx context.Context, questionID int64) error {
	return RegisterNoticeByQuestionID(ctx, r.DB, questionID)
}

// DeleteByQuestionID は質問 ID に紐づく通知を削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - questionID int64: 質問 ID
//
// return:
//   - int64: 削除した通知 ID（見つからない場合は -1）
//   - error: DB エラー
func (r *NoticeRepository) DeleteByQuestionID(ctx context.Context, questionID int64) (int64, error) {
	return DeleteNoticeByQuestionID(ctx, r.DB, questionID)
}
