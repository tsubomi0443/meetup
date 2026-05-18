package auth

import (
	"context"
	"errors"

	"meetup/internal/domains/entity"
	domainuser "meetup/internal/domains/user"
	"meetup/internal/ports"

	"gorm.io/gorm"
)

// UseCase は認証（ログイン）ユースケースを表す。
type UseCase struct {
	users    domainuser.Repository
	password ports.PasswordHasher
}

// NewUseCase は認証ユースケースを生成する。
//
// args:
//   - users domainuser.Repository: ユーザーリポジトリ
//   - password ports.PasswordHasher: パスワードハッシュ・照合
//
// return:
//   - *UseCase: 生成したユースケース
func NewUseCase(users domainuser.Repository, password ports.PasswordHasher) *UseCase {
	return &UseCase{users: users, password: password}
}

// LoginResult はログイン成功時に返すユーザ情報と保存済みパスワードハッシュを保持する。
type LoginResult struct {
	User         entity.User
	StoredPass   string
}

// Login はメールアドレスと平文パスワードでログイン認証を行う。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - email string: ログイン用メールアドレス
//   - plainPassword string: 平文パスワード
//
// return:
//   - LoginResult: 認証成功時のユーザ情報と保存済みハッシュ
//   - error: 認証失敗・DB エラー（未登録は gorm.ErrRecordNotFound）
func (u *UseCase) Login(ctx context.Context, email, plainPassword string) (LoginResult, error) {
	stored, err := u.users.GetPasswordByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResult{}, gorm.ErrRecordNotFound
		}
		return LoginResult{}, err
	}
	ok, err := u.password.VerifyPassword(stored, plainPassword)
	if err != nil {
		return LoginResult{}, err
	}
	if !ok {
		return LoginResult{}, gorm.ErrRecordNotFound
	}
	user, err := u.users.GetUserInfo(ctx, email, stored, "Role")
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{User: user, StoredPass: stored}, nil
}
