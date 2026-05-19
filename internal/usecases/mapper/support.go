package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// サポート（entity.Support）

// =====================
// SupportFromEntity は entity.Support を dto.SupportForm に変換する。
//
// args:
//   - e entity.Support: 変換元エンティティ
//
// return:
//   - dto.SupportForm: サポートフォーム DTO
func SupportFromEntity(e entity.Support) dto.SupportForm {
	f := dto.SupportForm{
		ID:              e.ID,
		UserID:          strconv.FormatInt(e.UserID, 10),
		SupportStatusID: strconv.FormatInt(e.SupportStatusID, 10),
		CreatedAt:       timeToISO(e.CreatedAt),
		UpdatedAt:       timeToISO(e.UpdatedAt),
		DeletedAt:       deletedAtToISO(e.DeletedAt),
	}
	if e.User.ID != 0 {
		u := UserFromEntity(e.User)
		f.User = &u
	}
	if e.SupportStatus.ID != 0 {
		ss := supportStatusFromEntityShallow(e.SupportStatus)
		f.SupportStatus = &ss
	}
	return f
}

// supportStatusFromEntityShallow は entity.SupportStatus を Supports なしの dto.SupportStatusForm に変換する。
//
// args:
//   - e entity.SupportStatus: 変換元エンティティ
//
// return:
//   - dto.SupportStatusForm: サポートステータスフォーム DTO（Supports なし）
func supportStatusFromEntityShallow(e entity.SupportStatus) dto.SupportStatusForm {
	return dto.SupportStatusForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// SupportToEntity は dto.SupportForm を entity.Support に変換する。
//
// args:
//   - f dto.SupportForm: 変換元フォーム DTO
//
// return:
//   - entity.Support: サポートエンティティ
func SupportToEntity(f dto.SupportForm) entity.Support {
	e := entity.Support{
		ID:              f.ID,
		UserID:          f.UserIDInt64(),
		SupportStatusID: f.SupportStatusIDInt64(),
	}
	if f.User != nil {
		e.User = UserToEntity(*f.User)
	}
	if f.SupportStatus != nil {
		e.SupportStatus = SupportStatusToEntity(*f.SupportStatus)
	}
	return e
}
