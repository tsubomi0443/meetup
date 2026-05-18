package tag

import (
	"context"

	domaintag "meetup/internal/domains/tag"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

type UseCase struct {
	tags domaintag.Repository
}

func NewUseCase(tags domaintag.Repository) *UseCase {
	return &UseCase{tags: tags}
}

func (u *UseCase) GetAll(ctx context.Context) ([]dto.TagForm, error) {
	models, err := u.tags.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.TagFromEntities(models), nil
}

func (u *UseCase) Register(ctx context.Context, form dto.TagForm) (dto.TagForm, error) {
	model := mapper.TagToEntity(form)
	if err := u.tags.Register(ctx, &model); err != nil {
		return dto.TagForm{}, err
	}
	loaded, err := u.tags.GetByID(ctx, model.ID)
	if err != nil {
		return dto.TagForm{}, err
	}
	return mapper.TagFromEntity(loaded), nil
}

func (u *UseCase) Update(ctx context.Context, form dto.TagForm) error {
	model := mapper.TagToEntityNoRelations(form)
	_, err := u.tags.UpdateByID(ctx, model.ID, model, "Category", "TagManagers")
	return err
}

func (u *UseCase) DeleteByID(ctx context.Context, id int64) error {
	return u.tags.DeleteByID(ctx, id)
}
