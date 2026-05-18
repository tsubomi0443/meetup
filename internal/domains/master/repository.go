package master

import (
	"context"

	"meetup/internal/domains/entity"
)

type Repository interface {
	GetRoles(ctx context.Context) ([]entity.Role, error)
	GetCategories(ctx context.Context) ([]entity.Category, error)
	GetSupportStatuses(ctx context.Context) ([]entity.SupportStatus, error)
}
