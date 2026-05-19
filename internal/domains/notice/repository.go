// Package notice は通知ドメインのリポジトリ契約を定義する。
package notice

import (
	"context"

	"meetup/internal/domains/entity"
)

// Repository は通知の永続化操作を抽象化する。
type Repository interface {
	// GetAll は全通知を取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.Notice: 通知のスライス
	//   - error: 取得に失敗した場合のエラー
	GetAll(ctx context.Context) ([]entity.Notice, error)

	// GetByQuestionSilent は質問に紐づくサイレント通知を取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - question entity.Question: 対象質問
	//
	// return:
	//   - entity.Notice: 取得した通知
	//   - error: 取得に失敗した場合のエラー
	GetByQuestionSilent(ctx context.Context, question entity.Question) (entity.Notice, error)

	// RegisterByQuestionID は質問 ID に紐づく通知を登録する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - questionID int64: 質問 ID
	//
	// return:
	//   - error: 登録に失敗した場合のエラー
	RegisterByQuestionID(ctx context.Context, questionID int64) error

	// DeleteByQuestionID は質問 ID に紐づく通知を削除する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - questionID int64: 質問 ID
	//
	// return:
	//   - int64: 削除された行数
	//   - error: 削除に失敗した場合のエラー
	DeleteByQuestionID(ctx context.Context, questionID int64) (int64, error)
}
