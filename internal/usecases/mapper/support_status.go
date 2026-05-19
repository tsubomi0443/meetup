package mapper

import (
	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// サポートステータス（entity.SupportStatus）

// =====================
// SupportStatusFromEntity は entity.SupportStatus を dto.SupportStatusForm に変換する。
//
// args:
//   - e entity.SupportStatus: 変換元エンティティ
//
// return:
//   - dto.SupportStatusForm: サポートステータスフォーム DTO
func SupportStatusFromEntity(e entity.SupportStatus) dto.SupportStatusForm {
	f := dto.SupportStatusForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	for _, s := range e.Supports {
		f.Supports = append(f.Supports, SupportFromEntity(s))
	}
	return f
}

// SupportStatusToEntity は dto.SupportStatusForm を entity.SupportStatus に変換する。
//
// args:
//   - f dto.SupportStatusForm: 変換元フォーム DTO
//
// return:
//   - entity.SupportStatus: サポートステータスエンティティ
func SupportStatusToEntity(f dto.SupportStatusForm) entity.SupportStatus {
	e := entity.SupportStatus{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, sf := range f.Supports {
		e.Supports = append(e.Supports, SupportToEntity(sf))
	}
	return e
}
