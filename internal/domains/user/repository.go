package user

import (
	"context"

	"meetup/internal/domains/entity"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (entity.User, error)
	GetPasswordByEmail(ctx context.Context, email string) (string, error)
	GetUserInfo(ctx context.Context, email, pass string, preloads ...string) (entity.User, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	Register(ctx context.Context, model *entity.User) error
	UpdateByID(ctx context.Context, id int64, model entity.User, omit ...string) (int, error)
	DeleteByID(ctx context.Context, id int64) error
}
