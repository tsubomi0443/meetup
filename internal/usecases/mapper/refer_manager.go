package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// ReferManagerFromEntity は entity.ReferManager を dto.ReferManagerForm に変換する。
//
// args:
//   - e entity.ReferManager: 変換元エンティティ
//
// return:
//   - dto.ReferManagerForm: 参照リンク管理フォーム DTO
func ReferManagerFromEntity(e entity.ReferManager) dto.ReferManagerForm {
	f := dto.ReferManagerForm{
		ID:        e.ID,
		AnswerID:  strconv.FormatInt(e.AnswerID, 10),
		ReferID:   strconv.FormatInt(e.ReferID, 10),
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if e.Answer.ID != 0 {
		a := AnswerFromEntity(e.Answer)
		f.Answer = &a
	}
	if e.Refer.ID != 0 {
		r := ReferFromEntity(e.Refer)
		f.Refer = &r
	}
	return f
}

// ReferManagerToEntity は dto.ReferManagerForm を entity.ReferManager に変換する。
//
// args:
//   - f dto.ReferManagerForm: 変換元フォーム DTO
//
// return:
//   - entity.ReferManager: 参照リンク管理エンティティ
func ReferManagerToEntity(f dto.ReferManagerForm) entity.ReferManager {
	e := entity.ReferManager{
		ID:       f.ID,
		AnswerID: f.AnswerIDInt64(),
		ReferID:  f.ReferIDInt64(),
	}
	if f.Answer != nil {
		e.Answer = AnswerToEntity(*f.Answer)
	}
	if f.Refer != nil {
		e.Refer = ReferToEntity(*f.Refer)
	}
	return e
}
