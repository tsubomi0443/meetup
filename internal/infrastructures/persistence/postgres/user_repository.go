package postgres

import (
	"context"

	"meetup/internal/domains/entity"

	"gorm.io/gorm"
)

// UserRepository はユーザーの永続化を担う。
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository は UserRepository を生成する。
//
// args:
//   - db *gorm.DB: データベース接続
//
// return:
//   - *UserRepository: リポジトリ
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetByID は指定 ID のユーザーを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func (r *UserRepository) GetByID(ctx context.Context, id int64) (entity.User, error) {
	return GetUserByID(ctx, r.DB, id)
}

// GetPasswordByEmail はメールアドレスに紐づく保存パスワードハッシュを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - email string: メールアドレス
//
// return:
//   - string: パスワードハッシュ
//   - error: DB エラー
func (r *UserRepository) GetPasswordByEmail(ctx context.Context, email string) (string, error) {
	return GetUserPasswordByEmail(ctx, r.DB, email)
}

// GetUserInfo はメールとパスワードハッシュでユーザーを取得する（Preload 指定可）。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - email string: メールアドレス
//   - pass string: 保存済みパスワードハッシュ
//   - preloads ...string: GORM Preload 名
//
// return:
//   - entity.User: ユーザー
//   - error: DB エラー
func (r *UserRepository) GetUserInfo(ctx context.Context, email, pass string, preloads ...string) (entity.User, error) {
	return GetUserInfo(ctx, r.DB, email, pass, preloads...)
}

// GetUsers は管理者以外のユーザー一覧を取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []entity.User: ユーザー一覧
//   - error: DB エラー
func (r *UserRepository) GetUsers(ctx context.Context) ([]entity.User, error) {
	return GetUsers(ctx, r.DB)
}

// Register はユーザーを新規登録する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - model *entity.User: 登録するユーザー
//
// return:
//   - error: DB エラー
func (r *UserRepository) Register(ctx context.Context, model *entity.User) error {
	return Register(ctx, r.DB, model)
}

// UpdateByID は指定 ID のユーザーを更新する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//   - model entity.User: 更新内容
//   - omit ...string: 更新から除外する関連・列
//
// return:
//   - int: 更新行数
//   - error: DB エラー
func (r *UserRepository) UpdateByID(ctx context.Context, id int64, model entity.User, omit ...string) (int, error) {
	return UpdateByID(ctx, r.DB, id, model, omit...)
}

// DeleteByID は指定 ID のユーザーを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//
// return:
//   - error: DB エラー
func (r *UserRepository) DeleteByID(ctx context.Context, id int64) error {
	return DeleteUserByID(ctx, r.DB, id)
}
