package mapper

import (
	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// 参照リンク（entity.Refer）

// =====================
// referFromEntityShallow は entity.Refer を Answers なしの dto.ReferForm に変換する。
//
// args:
//   - e entity.Refer: 変換元エンティティ
//
// return:
//   - dto.ReferForm: 参照リンクフォーム DTO
func referFromEntityShallow(e entity.Refer) dto.ReferForm {
	return dto.ReferForm{
		ID:        e.ID,
		Title:     e.Title,
		URL:       e.URL,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// ReferFromEntity は entity.Refer を dto.ReferForm に変換する。
//
// args:
//   - e entity.Refer: 変換元エンティティ
//
// return:
//   - dto.ReferForm: 参照リンクフォーム DTO
func ReferFromEntity(e entity.Refer) dto.ReferForm {
	f := referFromEntityShallow(e)
	for _, rm := range e.ReferManagers {
		if rm.Answer.ID != 0 {
			f.Answers = append(f.Answers, AnswerFromEntity(rm.Answer))
		}
	}
	return f
}

// ReferToEntity は dto.ReferForm を entity.Refer に変換する。
//
// args:
//   - f dto.ReferForm: 変換元フォーム DTO
//
// return:
//   - entity.Refer: 参照リンクエンティティ
func ReferToEntity(f dto.ReferForm) entity.Refer {
	e := entity.Refer{
		ID:    f.ID,
		Title: f.Title,
		URL:   f.URL,
	}
	referID := f.ID
	for _, af := range f.Answers {
		rm := entity.ReferManager{
			ReferID:  referID,
			AnswerID: af.ID,
		}
		if af.ID != 0 {
			rm.Answer = entity.Answer{ID: af.ID}
		}
		e.ReferManagers = append(e.ReferManagers, rm)
	}
	return e
}
