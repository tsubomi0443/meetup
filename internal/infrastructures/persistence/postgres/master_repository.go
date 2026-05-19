package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// MasterRepository はマスタデータの永続化を担う。
type MasterRepository struct {
	DB *gorm.DB
}

// NewMasterRepository は MasterRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *MasterRepository: リポジトリ
func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{DB: db}
}

// GetRoles はロールマスタ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Role: ロール一覧
//   - error: DB エラー
func (r *MasterRepository) GetRoles(ctx context.Context) ([]entity.Role, error) {
	return GetMasterData[entity.Role](ctx, r.DB)
}

// GetCategories はカテゴリマスタ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Category: カテゴリ一覧
//   - error: DB エラー
func (r *MasterRepository) GetCategories(ctx context.Context) ([]entity.Category, error) {
	return GetMasterData[entity.Category](ctx, r.DB)
}

// GetSupportStatuses は支援ステータスマスタ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.SupportStatus: 支援ステータス一覧
//   - error: DB エラー
func (r *MasterRepository) GetSupportStatuses(ctx context.Context) ([]entity.SupportStatus, error) {
	return GetMasterData[entity.SupportStatus](ctx, r.DB)
}
