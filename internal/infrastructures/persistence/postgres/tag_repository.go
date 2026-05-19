package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// TagRepository はタグの永続化を担う。
type TagRepository struct {
	DB *gorm.DB
}

// NewTagRepository は TagRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *TagRepository: リポジトリ
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

// GetAll はタグ一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.Tag: タグ一覧
//   - error: DB エラー
func (r *TagRepository) GetAll(ctx context.Context) ([]entity.Tag, error) {
	return GetTags(ctx, r.DB)
}

// GetByID は指定 ID のタグを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//
// return:
//   - entity.Tag: タグ
//   - error: DB エラー
func (r *TagRepository) GetByID(ctx context.Context, id int64) (entity.Tag, error) {
	return GetTagByID(ctx, r.DB, id)
}

// Register はタグを新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - model *entity.Tag: 登録するタグ
//
// return:
//   - error: DB エラー
func (r *TagRepository) Register(ctx context.Context, model *entity.Tag) error {
	return Register(ctx, r.DB, model)
}

// UpdateByID は指定 ID のタグを更新する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//   - model entity.Tag: 更新内容
//   - omit ...string: 更新から除外する関連・列
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func (r *TagRepository) UpdateByID(ctx context.Context, id int64, model entity.Tag, omit ...string) (int, error) {
	return UpdateByID(ctx, r.DB, id, model, omit...)
}

// DeleteByID は指定 ID のタグを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: タグ ID
//
// return:
//   - error: DB エラー
func (r *TagRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteTagByID(ctx, r.DB, id)
}
