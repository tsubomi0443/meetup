package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// ユーザー（entity.User）

// =====================
// UserFromEntity は entity.User を dto.UserForm に変換する。Role があればネストして変換する。
//
// args:
//   - e entity.User: 変換元エンティティ
//
// return:
//   - dto.UserForm: ユーザーフォーム DTO
func UserFromEntity(e entity.User) dto.UserForm {
	f := dto.UserForm{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Memo:      e.Memo,
		RoleID:    strconv.FormatInt(e.RoleID, 10),
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if e.Role.ID != 0 {
		r := roleFromEntityShallow(e.Role)
		f.Role = &r
	}
	return f
}

// UserFromEntityNoRole は entity.User を dto.UserForm に変換する。entity.Role.Users 埋め込み時の循環参照を避けるため Role は含めない。
//
// args:
//   - e entity.User: 変換元エンティティ
//
// return:
//   - dto.UserForm: ユーザーフォーム DTO（Role ネストなし）
func UserFromEntityNoRole(e entity.User) dto.UserForm {
	return dto.UserForm{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Memo:      e.Memo,
		RoleID:    strconv.FormatInt(e.RoleID, 10),
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// UserToEntityNoRole は dto.UserForm を entity.User に変換する。Role の関連グラフは展開しない。
//
// args:
//   - f dto.UserForm: 変換元フォーム DTO
//
// return:
//   - entity.User: ユーザーエンティティ（Role は空構造体）
func UserToEntityNoRole(f dto.UserForm) entity.User {
	e := entity.User{
		ID:     f.ID,
		Name:   f.Name,
		Email:  f.Email,
		Memo:   f.Memo,
		RoleID: f.RoleIDInt64(),
	}
	if f.RoleID == "0" && f.Role != nil {
		e.RoleID = f.Role.ID
	}
	if f.Password != "" {
		e.Password = f.Password
	}
	e.Role = entity.Role{}
	return e
}

// UserToEntity は dto.UserForm を entity.User に変換する。Role があればネストして変換する。
//
// args:
//   - f dto.UserForm: 変換元フォーム DTO
//
// return:
//   - entity.User: ユーザーエンティティ
func UserToEntity(f dto.UserForm) entity.User {
	e := entity.User{
		ID:     f.ID,
		Name:   f.Name,
		Email:  f.Email,
		Memo:   f.Memo,
		RoleID: f.RoleIDInt64(),
	}
	if f.RoleID == "0" && f.Role != nil {
		e.RoleID = f.Role.ID
	}
	if f.Password != "" {
		e.Password = f.Password
	}
	if f.Role != nil {
		e.Role = RoleToEntity(*f.Role)
	}
	return e
}

// UserFormsFromEntities は entity.User のスライスを dto.UserForm のスライスに一括変換する。
//
// args:
//   - users []entity.User: 変換元エンティティ一覧
//
// return:
//   - []dto.UserForm: ユーザーフォーム DTO の一覧
func UserFormsFromEntities(users []entity.User) []dto.UserForm {
	out := make([]dto.UserForm, len(users))
	for i := range users {
		out[i] = UserFromEntity(users[i])
	}
	return out
}
