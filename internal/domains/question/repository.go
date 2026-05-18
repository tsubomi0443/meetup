package question

import (
	"context"

	"meetup/internal/domains/entity"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (entity.Question, error)
	GetAll(ctx context.Context) ([]entity.Question, error)
	Register(ctx context.Context, model *entity.Question) error
	UpdateInTransaction(ctx context.Context, q entity.Question) error
	DeleteByID(ctx context.Context, id int64) error
}
