package question

import (
	"context"

	"meetup/internal/domains/entity"
	domainquestion "meetup/internal/domains/question"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
)

type UseCase struct {
	questions domainquestion.Repository
}

func NewUseCase(questions domainquestion.Repository) *UseCase {
	return &UseCase{questions: questions}
}

func (u *UseCase) GetAll(ctx context.Context) ([]dto.QuestionForm, error) {
	models, err := u.questions.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.QuestionFromEntities(models), nil
}

func (u *UseCase) GetByID(ctx context.Context, id int64) (dto.QuestionForm, error) {
	model, err := u.questions.GetByID(ctx, id)
	if err != nil {
		return dto.QuestionForm{}, err
	}
	return mapper.QuestionFromEntity(model), nil
}

func (u *UseCase) Register(ctx context.Context, form dto.QuestionForm) (dto.QuestionForm, error) {
	data := mapper.QuestionToEntity(form)
	if err := u.questions.Register(ctx, &data); err != nil {
		return dto.QuestionForm{}, err
	}
	created, err := u.questions.GetByID(ctx, data.ID)
	if err != nil {
		return dto.QuestionForm{}, err
	}
	return mapper.QuestionFromEntity(created), nil
}

func (u *UseCase) Update(ctx context.Context, form dto.QuestionForm, actorUserID int64, hasActor bool) (entity.Question, dto.QuestionForm, error) {
	NormalizeQuestionFormClearSupportWhenUnassigned(&form)
	if hasActor {
		NormalizeQuestionFormAssignSupportUserWhenInProgress(&form, actorUserID)
	}
	updatedModel := mapper.QuestionToEntity(form)
	if err := u.questions.UpdateInTransaction(ctx, updatedModel); err != nil {
		return entity.Question{}, dto.QuestionForm{}, err
	}
	loaded, err := u.questions.GetByID(ctx, updatedModel.ID)
	if err != nil {
		return entity.Question{}, dto.QuestionForm{}, err
	}
	return updatedModel, mapper.QuestionFromEntity(loaded), nil
}

func (u *UseCase) DeleteByID(ctx context.Context, id int64) error {
	return u.questions.DeleteByID(ctx, id)
}
