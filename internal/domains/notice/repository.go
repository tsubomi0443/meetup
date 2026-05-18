package notice

import (
	"context"

	"meetup/internal/domains/entity"
)

type Repository interface {
	GetAll(ctx context.Context) ([]entity.Notice, error)
	GetByQuestionSilent(ctx context.Context, question entity.Question) (entity.Notice, error)
	RegisterByQuestionID(ctx context.Context, questionID int64) error
	DeleteByQuestionID(ctx context.Context, questionID int64) (int64, error)
}
