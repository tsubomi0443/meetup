package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// TagManagerFromEntity は entity.TagManager を dto.TagManagerForm に変換する。
//
// args:
//   - e entity.TagManager: 変換元エンティティ
//
// return:
//   - dto.TagManagerForm: タグ管理フォーム DTO
func TagManagerFromEntity(e entity.TagManager) dto.TagManagerForm {
	f := dto.TagManagerForm{
		ID:         e.ID,
		TagID:      strconv.FormatInt(e.TagID, 10),
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.Tag.ID != 0 {
		t := TagFromEntity(e.Tag)
		f.Tag = &t
	}
	if e.Question.ID != 0 {
		q := QuestionFromEntity(e.Question)
		f.Question = &q
	}
	return f
}

// TagManagerToEntity は dto.TagManagerForm を entity.TagManager に変換する。
//
// args:
//   - f dto.TagManagerForm: 変換元フォーム DTO
//
// return:
//   - entity.TagManager: タグ管理エンティティ
func TagManagerToEntity(f dto.TagManagerForm) entity.TagManager {
	e := entity.TagManager{
		ID:         f.ID,
		TagID:      f.TagIDInt64(),
		QuestionID: f.QuestionIDInt64(),
	}
	if f.Tag != nil {
		e.Tag = TagToEntity(*f.Tag)
	}
	if f.Question != nil {
		e.Question = QuestionToEntity(*f.Question)
	}
	return e
}
