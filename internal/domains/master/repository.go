// Package master はマスタデータのリポジトリ契約を定義する。
package master

import (
	"context"

	"meetup/internal/domains/entity"
)

// Repository はマスタデータの取得操作を抽象化する。
type Repository interface {
	// GetRoles は全ロールを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.Role: ロールのスライス
	//   - error: 取得に失敗した場合のエラー
	GetRoles(ctx context.Context) ([]entity.Role, error)

	// GetCategories は全カテゴリを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.Category: カテゴリのスライス
	//   - error: 取得に失敗した場合のエラー
	GetCategories(ctx context.Context) ([]entity.Category, error)

	// GetSupportStatuses は全サポートステータスを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.SupportStatus: サポートステータスのスライス
	//   - error: 取得に失敗した場合のエラー
	GetSupportStatuses(ctx context.Context) ([]entity.SupportStatus, error)
}
