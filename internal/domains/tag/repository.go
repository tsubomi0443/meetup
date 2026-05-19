// Package tag はタグドメインのリポジトリ契約を定義する。
package tag

import (
	"context"

	"meetup/internal/domains/entity"
)

// Repository はタグの永続化操作を抽象化する。
type Repository interface {
	// GetAll は全タグを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.Tag: タグのスライス
	//   - error: 取得に失敗した場合のエラー
	GetAll(ctx context.Context) ([]entity.Tag, error)

	// GetByID は指定 ID のタグを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: タグ ID
	//
	// return:
	//   - entity.Tag: 取得したタグ
	//   - error: 取得に失敗した場合のエラー
	GetByID(ctx context.Context, id int64) (entity.Tag, error)

	// Register は新規タグを登録する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - model *entity.Tag: 登録するタグ（ID は保存後に設定される）
	//
	// return:
	//   - error: 登録に失敗した場合のエラー
	Register(ctx context.Context, model *entity.Tag) error

	// UpdateByID は指定 ID のタグを更新する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: タグ ID
	//   - model entity.Tag: 更新内容
	//   - omit ...string: 更新から除外するカラム名
	//
	// return:
	//   - int: 更新された行数
	//   - error: 更新に失敗した場合のエラー
	UpdateByID(ctx context.Context, id int64, model entity.Tag, omit ...string) (int, error)

	// DeleteByID は指定 ID のタグを論理削除する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: タグ ID
	//
	// return:
	//   - error: 削除に失敗した場合のエラー
	DeleteByID(ctx context.Context, id int64) error
}
