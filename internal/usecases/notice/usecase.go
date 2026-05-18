package notice

import (
	"context"

	domainnotice "meetup/internal/domains/notice"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

type UseCase struct {
	notices domainnotice.Repository
}

func NewUseCase(notices domainnotice.Repository) *UseCase {
	return &UseCase{notices: notices}
}

func (u *UseCase) GetAll(ctx context.Context) ([]dto.NoticeForm, error) {
	models, err := u.notices.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.NoticeFromEntities(models), nil
}
