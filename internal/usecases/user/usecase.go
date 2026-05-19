package user

import (
	"context"
	"strings"

	domainuser "meetup/internal/domains/user"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
	"meetup/internal/ports"
)

// UseCase はユーザー管理ユースケースを表す。
type UseCase struct {
	users    domainuser.Repository
	password ports.PasswordHasher
}

// NewUseCase はユーザー管理ユースケースを生成する。
//
// args:
//   - users domainuser.Repository: ユーザーリポジトリ
//   - password ports.PasswordHasher: パスワードハッシュ・暗号化
//
// return:
//   - *UseCase: 生成したユースケース
func NewUseCase(users domainuser.Repository, password ports.PasswordHasher) *UseCase {
	return &UseCase{users: users, password: password}
}

// GetAll は全ユーザーを DTO 一覧として取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//
// return:
//   - []dto.UserForm: ユーザーフォームの一覧
//   - error: 取得エラー
func (u *UseCase) GetAll(ctx context.Context) ([]dto.UserForm, error) {
	models, err := u.users.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.UserFormsFromEntities(models), nil
}

// GetByID は指定 ID のユーザーを取得する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//
// return:
//   - dto.UserForm: ユーザーフォーム
//   - error: 取得エラー
func (u *UseCase) GetByID(ctx context.Context, id int64) (dto.UserForm, error) {
	model, err := u.users.GetByID(ctx, id)
	if err != nil {
		return dto.UserForm{}, err
	}
	return mapper.UserFromEntity(model), nil
}

// Register は新規ユーザーを登録する。パスワードが指定されていれば Argon2 でハッシュ化する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - form dto.UserForm: 登録内容
//
// return:
//   - dto.UserForm: 登録後のユーザーフォーム（パスワードは空）
//   - error: 登録・取得エラー
func (u *UseCase) Register(ctx context.Context, form dto.UserForm) (dto.UserForm, error) {
	if strings.TrimSpace(form.Password) != "" {
		enc, err := u.password.EncryptPasswordByArgon2Encode(form.Password)
		if err != nil {
			return dto.UserForm{}, err
		}
		form.Password = enc
	}
	data := mapper.UserToEntityNoRole(form)
	if err := u.users.Register(ctx, &data); err != nil {
		return dto.UserForm{}, err
	}
	created, err := u.users.GetByID(ctx, data.ID)
	if err != nil {
		return dto.UserForm{}, err
	}
	created.Password = ""
	return mapper.UserFromEntity(created), nil
}

// Update は既存ユーザーを更新する。パスワードが指定されていれば Argon2 でハッシュ化する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - form dto.UserForm: 更新内容
//
// return:
//   - error: 更新エラー
func (u *UseCase) Update(ctx context.Context, form dto.UserForm) error {
	if strings.TrimSpace(form.Password) != "" {
		enc, err := u.password.EncryptPasswordByArgon2Encode(form.Password)
		if err != nil {
			return err
		}
		form.Password = enc
	}
	updatedModel := mapper.UserToEntityNoRole(form)
	_, err := u.users.UpdateByID(ctx, updatedModel.ID, updatedModel, "Role")
	return err
}

// DeleteByID は指定 ID のユーザーを削除する。
//
// args:
//   - ctx context.Context: リクエストコンテキスト
//   - id int64: ユーザー ID
//
// return:
//   - error: 削除エラー
func (u *UseCase) DeleteByID(ctx context.Context, id int64) error {
	return u.users.DeleteByID(ctx, id)
}
