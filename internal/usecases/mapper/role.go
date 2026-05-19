// Package mapper はドメインエンティティと dto フォーム DTO の相互変換を提供する。
package mapper

import (
	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// ロール（entity.Role）

// =====================
// RoleFromEntity は entity.Role を dto.RoleForm に変換する。Users があればネストして変換する。
//
// args:
//   - e entity.Role: 変換元エンティティ
//
// return:
//   - dto.RoleForm: ロールフォーム DTO
func RoleFromEntity(e entity.Role) dto.RoleForm {
	f := dto.RoleForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if len(e.Users) > 0 {
		f.Users = make([]dto.UserForm, len(e.Users))
		for i := range e.Users {
			f.Users[i] = UserFromEntityNoRole(e.Users[i])
		}
	}
	return f
}

// roleFromEntityShallow は entity.Role を関連 Users なしの dto.RoleForm に変換する（循環参照回避用）。
//
// args:
//   - e entity.Role: 変換元エンティティ
//
// return:
//   - dto.RoleForm: ロールフォーム DTO（Users なし）
func roleFromEntityShallow(e entity.Role) dto.RoleForm {
	return dto.RoleForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// RoleToEntity は dto.RoleForm を entity.Role に変換する。
//
// args:
//   - f dto.RoleForm: 変換元フォーム DTO
//
// return:
//   - entity.Role: ロールエンティティ
func RoleToEntity(f dto.RoleForm) entity.Role {
	e := entity.Role{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, uf := range f.Users {
		e.Users = append(e.Users, UserToEntity(uf))
	}
	return e
}
