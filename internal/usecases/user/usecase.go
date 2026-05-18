package user

import (
	"context"
	"strings"

	domainuser "meetup/internal/domains/user"
	"meetup/internal/usecases/mapper"
	"meetup/internal/usecases/dto"
	"meetup/internal/ports"
)

type UseCase struct {
	users    domainuser.Repository
	password ports.PasswordHasher
}

func NewUseCase(users domainuser.Repository, password ports.PasswordHasher) *UseCase {
	return &UseCase{users: users, password: password}
}

func (u *UseCase) GetAll(ctx context.Context) ([]dto.UserForm, error) {
	models, err := u.users.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.UserFormsFromEntities(models), nil
}

func (u *UseCase) GetByID(ctx context.Context, id int64) (dto.UserForm, error) {
	model, err := u.users.GetByID(ctx, id)
	if err != nil {
		return dto.UserForm{}, err
	}
	return mapper.UserFromEntity(model), nil
}

func (u *UseCase) Register(ctx context.Context, form dto.UserForm) (dto.UserForm, error) {
	if strings.TrimSpace(form.Password) != "" {
		enc, err := u.password.EncryptPasswordByArgon2Encode(form.Password)
		if err != nil {
			return dto.UserForm{}, err
		}
		form.Password = enc
	}
	data := mapper.UserToEntityNoRole(form)
	if err := u.users.Register(ctx, &data); err != nil {
		return dto.UserForm{}, err
	}
	created, err := u.users.GetByID(ctx, data.ID)
	if err != nil {
		return dto.UserForm{}, err
	}
	created.Password = ""
	return mapper.UserFromEntity(created), nil
}

func (u *UseCase) Update(ctx context.Context, form dto.UserForm) error {
	if strings.TrimSpace(form.Password) != "" {
		enc, err := u.password.EncryptPasswordByArgon2Encode(form.Password)
		if err != nil {
			return err
		}
		form.Password = enc
	}
	updatedModel := mapper.UserToEntityNoRole(form)
	_, err := u.users.UpdateByID(ctx, updatedModel.ID, updatedModel, "Role")
	return err
}

func (u *UseCase) DeleteByID(ctx context.Context, id int64) error {
	return u.users.DeleteByID(ctx, id)
}
