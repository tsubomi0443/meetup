package tag

import (
	"context"

	"meetup/internal/domains/entity"
)

type Repository interface {
	GetAll(ctx context.Context) ([]entity.Tag, error)
	GetByID(ctx context.Context, id int64) (entity.Tag, error)
	Register(ctx context.Context, model *entity.Tag) error
	UpdateByID(ctx context.Context, id int64, model entity.Tag, omit ...string) (int, error)
	DeleteByID(ctx context.Context, id int64) error
}
