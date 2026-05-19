// Package user はユーザードメインのリポジトリ契約を定義する。
package user

import (
	"context"

	"meetup/internal/domains/entity"
)

// Repository はユーザーの永続化操作を抽象化する。
type Repository interface {
	// GetByID は指定 ID のユーザーを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: ユーザー ID
	//
	// return:
	//   - entity.User: 取得したユーザー
	//   - error: 取得に失敗した場合のエラー
	GetByID(ctx context.Context, id int64) (entity.User, error)

	// GetPasswordByEmail はメールアドレスに紐づくパスワードハッシュを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - email string: メールアドレス
	//
	// return:
	//   - string: パスワードハッシュ
	//   - error: 取得に失敗した場合のエラー
	GetPasswordByEmail(ctx context.Context, email string) (string, error)

	// GetUserInfo はメールアドレスとパスワードでユーザーを認証付き取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - email string: メールアドレス
	//   - pass string: パスワード（平文）
	//   - preloads ...string: GORM でプリロードする関連名
	//
	// return:
	//   - entity.User: 取得したユーザー
	//   - error: 認証または取得に失敗した場合のエラー
	GetUserInfo(ctx context.Context, email, pass string, preloads ...string) (entity.User, error)

	// GetUsers は全ユーザーを取得する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//
	// return:
	//   - []entity.User: ユーザーのスライス
	//   - error: 取得に失敗した場合のエラー
	GetUsers(ctx context.Context) ([]entity.User, error)

	// Register は新規ユーザーを登録する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - model *entity.User: 登録するユーザー（ID は保存後に設定される）
	//
	// return:
	//   - error: 登録に失敗した場合のエラー
	Register(ctx context.Context, model *entity.User) error

	// UpdateByID は指定 ID のユーザーを更新する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: ユーザー ID
	//   - model entity.User: 更新内容
	//   - omit ...string: 更新から除外するカラム名
	//
	// return:
	//   - int: 更新された行数
	//   - error: 更新に失敗した場合のエラー
	UpdateByID(ctx context.Context, id int64, model entity.User, omit ...string) (int, error)

	// DeleteByID は指定 ID のユーザーを論理削除する。
	//
	// args:
	//   - ctx context.Context: リクエストコンテキスト
	//   - id int64: ユーザー ID
	//
	// return:
	//   - error: 削除に失敗した場合のエラー
	DeleteByID(ctx context.Context, id int64) error
}
