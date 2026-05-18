package auth

import (
	"context"
	"errors"

	"meetup/internal/domains/entity"
	domainuser "meetup/internal/domains/user"
	"meetup/internal/ports"

	"gorm.io/gorm"
)

type UseCase struct {
	users    domainuser.Repository
	password ports.PasswordHasher
}

func NewUseCase(users domainuser.Repository, password ports.PasswordHasher) *UseCase {
	return &UseCase{users: users, password: password}
}

// LoginResult holds user data after successful credential check.
type LoginResult struct {
	User         entity.User
	StoredPass   string
}

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
