// Package question は質問ドメインのリポジトリ契約を定義する。
package question

import (
	"context"

	"meetup/internal/domains/entity"
)

// Repository は質問の永続化操作を抽象化する。
type Repository interface {
	// GetByID は指定 ID の質問を取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: 質問 ID
	//
	// return:
	//   - entity.Question: 取得した質問
	//   - error: 取得に失敗した場合のエラー
	GetByID(ctx context.Context, id int64) (entity.Question, error)

	// GetAll は全質問を取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.Question: 質問のスライス
	//   - error: 取得に失敗した場合のエラー
	GetAll(ctx context.Context) ([]entity.Question, error)

	// Register は新規質問を登録する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - model *entity.Question: 登録する質問（ID は保存後に設定される）
	//
	// return:
	//   - error: 登録に失敗した場合のエラー
	Register(ctx context.Context, model *entity.Question) error

	// UpdateInTransaction はトランザクション内で質問を更新する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - q entity.Question: 更新する質問
	//
	// return:
	//   - error: 更新に失敗した場合のエラー
	UpdateInTransaction(ctx context.Context, q entity.Question) error

	// DeleteByID は指定 ID の質問を論理削除する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: 質問 ID
	//
	// return:
	//   - error: 削除に失敗した場合のエラー
	DeleteByID(ctx context.Context, id int64) error
}
