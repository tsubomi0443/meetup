package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// GetUserByID は指定 ID のユーザーをロール込みで取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: ユーザー ID
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func GetUserByID(ctx context.Context, db *gorm.DB, id int64) (model entity.User, err error) {
	model, err = gorm.G[entity.User](db).Where("id = ?", id).Select("id, name, email, memo, role_id").
		Preload("Role", commonPreloadBuilder()).
		First(ctx)
	return
}

// GetUserPasswordByEmail はメールアドレスに紐づく保存パスワードハッシュを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - email string: メールアドレス
//
// return:
//   - string: パスワードハッシュ
//   - error: DB エラー
func GetUserPasswordByEmail(ctx context.Context, db *gorm.DB, email string) (password string, err error) {
	u, err := gorm.G[entity.User](db).Where("email = ?", email).Select("id, password").First(ctx)
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

// GetUserInfo はメールとパスワードハッシュでユーザーを取得する（Preload 指定可）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - email string: メールアドレス
//   - pass string: 保存済みパスワードハッシュ
//   - preloads ...string: GORM Preload 名
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func GetUserInfo(ctx context.Context, db *gorm.DB, email, pass string, preloads ...string) (model entity.User, err error) {
	chain := gorm.G[entity.User](db).Where("email = ? AND password = ?", email, pass)
	for _, preload := range preloads {
		chain = chain.Preload(preload, commonPreloadBuilder())
	}
	model, err = chain.First(ctx)
	return
}

// GetUsers は管理者（role_id=1）以外のユーザー一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//
// return:
//   - []entity.User: ユーザー一覧
//   - error: DB エラー
func GetUsers(ctx context.Context, db *gorm.DB) (models []entity.User, err error) {
	models, err = gorm.G[entity.User](db).
		Where("role_id <> ?", 1).
		Preload("Role", commonPreloadBuilder()).
		Not("role_id = 1").
		Select("id, name, email, memo, role_id").
		Order("id").
		Find(ctx)
	return
}

// DeleteUserByID は指定 ID のユーザーを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - db *gorm.DB: データベース接続
//   - id int64: ユーザー ID
//
// return:
//   - error: DB エラー
func DeleteUserByID(ctx context.Context, db *gorm.DB, id int64) error {
	if _, err := gorm.G[entity.User](db).
		Preload("Role", commonPreloadBuilder()).
		Where("id = ?", id).
		Limit(1).
		Delete(ctx); err != nil {
		return err
	}
	return nil
}
