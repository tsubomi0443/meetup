package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// GetMasterData は entity.Role / entity.Category / entity.SupportStatus などマスタテーブルの一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []T: マスタ行の一覧
//   - error: DB エラー
func GetMasterData[T entity.Role | entity.Category | entity.SupportStatus](ctx context.Context, db *gorm.DB) (models []T, err error) {
	models, err = gorm.G[T](db).Find(ctx)
	return
}

// Register は任意型のモデルを新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - model T: 登録するモデル
//   - preloads ...string: 作成後 Preload 名（未使用の場合あり）
//
// return:
//   - error: DB エラー
func Register[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) error {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Create(ctx, &model)
}

// Updates はモデルを一括更新する（WHERE なしの Updates）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - model T: 更新内容
//   - preloads ...string: Preload 名（未使用の場合あり）
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func Updates[T any](ctx context.Context, db *gorm.DB, model T, preloads ...string) (int, error) {
	var v = gorm.G[T](db)
	for _, preload := range preloads {
		v.Preload(preload, nil)
	}
	return v.Updates(ctx, model)
}

// UpdateByID は主キー id を WHERE に固定して単一モデルを更新する。
// gorm.G[T](db) の Updates を利用し、トランザクションは張らない。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: 主キー
//   - model T: 更新内容
//   - omit ...string: Omit に渡す関連名・列名
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func UpdateByID[T any](ctx context.Context, db *gorm.DB, id int64, model T, omit ...string) (int, error) {
	return gorm.G[T](db.WithContext(ctx)).
		Omit(omit...).
		Where("id = ?", id).
		Updates(ctx, model)
}

// UpdateInTransaction は単一トランザクション内で gorm.Updates を実行する。
// omit に関連名を渡し、中間テーブル向け関連を更新対象から外す。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - model T: 更新内容
//   - omit ...string: 更新から除外する関連名
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func UpdateInTransaction[T any](ctx context.Context, db *gorm.DB, model T, omit ...string) (rowsAffected int, err error) {
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := model
		res := tx.Omit(omit...).Model(&m).Updates(&m)
		rowsAffected = int(res.RowsAffected)
		return res.Error
	})
	return
}

// GetMaxByColumn は指定列の MAX 値を返す（無効・未存在時は -1）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - columnName string: 集計列名
//
// return:
//   - int64: 最大値
func GetMaxByColumn[T any](ctx context.Context, db *gorm.DB, columnName string) int64 {
	var max sql.NullInt64
	err := db.WithContext(ctx).Model(new(T)).
		Select(fmt.Sprintf("MAX(%s)", columnName)).
		Take(&max).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return -1
	}

	if !max.Valid {
		return -1
	}
	return max.Int64

}
