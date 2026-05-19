package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// GetTags はタグ一覧をカテゴリ込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.Tag: タグ一覧
//   - error: DB エラー
func GetTags(ctx context.Context, db *gorm.DB) (models []entity.Tag, err error) {
	models, err = gorm.G[entity.Tag](db).
		Preload("Category", commonPreloadBuilder()).
		Order("id").
		Find(ctx)
	return
}

// GetTagByID は指定 ID のタグをカテゴリ込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: タグ ID
//
// return:
//   - entity.Tag: タグ
//   - error: DB エラー
func GetTagByID(ctx context.Context, db *gorm.DB, id int64) (models entity.Tag, err error) {
	models, err = gorm.G[entity.Tag](db).
		Preload("Category", commonPreloadBuilder()).
		Where("id = ?", id).
		First(ctx)
	return
}

// DeleteTagByID は指定 ID のタグを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: タグ ID
//
// return:
//   - error: DB エラー
func DeleteTagByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[entity.Tag](db).Where("id = ?", id).Limit(1).Delete(ctx); err != nil {
		return err
	}
	return nil
}
