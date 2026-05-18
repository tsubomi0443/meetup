package master

import (
	"context"

	domainmaster "meetup/internal/domains/master"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

type UseCase struct {
	master domainmaster.Repository
}

func NewUseCase(master domainmaster.Repository) *UseCase {
	return &UseCase{master: master}
}

func (u *UseCase) ListRoles(ctx context.Context) ([]dto.RoleForm, error) {
	roles, err := u.master.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	roleForms := make([]dto.RoleForm, 0, len(roles))
	for i := range roles {
		roleForms = append(roleForms, mapper.RoleFromEntity(roles[i]))
	}
	return roleForms, nil
}

func (u *UseCase) ListCategories(ctx context.Context) ([]dto.CategoryForm, error) {
	categories, err := u.master.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	categoryForms := make([]dto.CategoryForm, 0, len(categories))
	for i := range categories {
		categoryForms = append(categoryForms, mapper.CategoryFromEntity(categories[i]))
	}
	return categoryForms, nil
}

func (u *UseCase) ListSupportStatuses(ctx context.Context) ([]dto.SupportStatusForm, error) {
	statuses, err := u.master.GetSupportStatuses(ctx)
	if err != nil {
		return nil, err
	}
	supportStatusForms := make([]dto.SupportStatusForm, 0, len(statuses))
	for i := range statuses {
		supportStatusForms = append(supportStatusForms, mapper.SupportStatusFromEntity(statuses[i]))
	}
	return supportStatusForms, nil
}
