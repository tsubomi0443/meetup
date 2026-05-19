package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// noticeTypeFromEntityShallow は entity.NoticeType を Notices なしの dto.NoticeTypeForm に変換する。
//
// args:
//   - e entity.NoticeType: 変換元エンティティ
//
// return:
//   - dto.NoticeTypeForm: 通知種別フォーム DTO
func noticeTypeFromEntityShallow(e entity.NoticeType) dto.NoticeTypeForm {
	return dto.NoticeTypeForm{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
}

// NoticeTypeFromEntity は entity.NoticeType を dto.NoticeTypeForm に変換する。
//
// args:
//   - e entity.NoticeType: 変換元エンティティ
//
// return:
//   - dto.NoticeTypeForm: 通知種別フォーム DTO
func NoticeTypeFromEntity(e entity.NoticeType) dto.NoticeTypeForm {
	f := noticeTypeFromEntityShallow(e)
	for _, n := range e.Notices {
		f.Notices = append(f.Notices, NoticeFromEntity(n))
	}
	return f
}

// NoticeTypeToEntity は dto.NoticeTypeForm を entity.NoticeType に変換する。
//
// args:
//   - f dto.NoticeTypeForm: 変換元フォーム DTO
//
// return:
//   - entity.NoticeType: 通知種別エンティティ
func NoticeTypeToEntity(f dto.NoticeTypeForm) entity.NoticeType {
	e := entity.NoticeType{
		ID:   f.ID,
		Name: f.Name,
	}
	for _, nf := range f.Notices {
		e.Notices = append(e.Notices, NoticeToEntity(nf))
	}
	return e
}

// NoticeFromEntity は entity.Notice を dto.NoticeForm に変換する。
//
// args:
//   - e entity.Notice: 変換元エンティティ
//
// return:
//   - dto.NoticeForm: 通知フォーム DTO
func NoticeFromEntity(e entity.Notice) dto.NoticeForm {
	var questionID *string
	if e.QuestionID != nil {
		s := strconv.FormatInt(*e.QuestionID, 10)
		questionID = &s
	}
	f := dto.NoticeForm{
		ID:         e.ID,
		TypeID:     e.TypeID,
		QuestionID: questionID,
		Content:    e.Content,
		DisplayDue: timePtrToISO(e.DisplayDue),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
	if e.NoticeType.ID != 0 {
		nt := noticeTypeFromEntityShallow(e.NoticeType)
		f.NoticeType = &nt
	}
	if e.Question != nil && e.Question.ID != 0 {
		// due := e.Question.Due.Format(time.DateTime)
		// const tag = TagFromEntity(e.Question.TagManagers)
		// qf := dto.QuestionForm{ID: e.Question.ID, Title: e.Question.Title, Content: e.Question.Content, Due: &due, Tags: tag}
		qf := QuestionFromEntity(*e.Question)
		f.Question = &qf
	}
	return f
}

// NoticeFromEntities は entity.Notice のスライスを dto.NoticeForm のスライスに一括変換する。
//
// args:
//   - ns []entity.Notice: 変換元エンティティ一覧
//
// return:
//   - []dto.NoticeForm: 通知フォーム DTO の一覧
func NoticeFromEntities(ns []entity.Notice) (forms []dto.NoticeForm) {
	for _, n := range ns {
		forms = append(forms, NoticeFromEntity(n))
	}
	return forms
}

// NoticeToEntity は dto.NoticeForm を entity.Notice に変換する。
//
// args:
//   - f dto.NoticeForm: 変換元フォーム DTO
//
// return:
//   - entity.Notice: 通知エンティティ
func NoticeToEntity(f dto.NoticeForm) entity.Notice {
	var questionID *int64
	if v := f.QuestionIDInt64(); v >= 0 {
		questionID = &v
	}
	e := entity.Notice{
		ID:         f.ID,
		TypeID:     f.TypeID,
		QuestionID: questionID,
		Content:    f.Content,
		DisplayDue: isoToTimePtr(f.DisplayDue),
	}
	if f.NoticeType != nil {
		e.NoticeType = NoticeTypeToEntity(*f.NoticeType)
	}
	if f.Question != nil {
		q := QuestionToEntity(*f.Question)
		e.Question = &q
	}
	return e
}
