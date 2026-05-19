package mapper

import (
	"strconv"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// 送信者・送信者トーク（entity.Sender / entity.SenderTalk）

// =====================
// SenderFromEntity は entity.Sender を dto.SenderForm に変換する。
//
// args:
//   - e entity.Sender: 変換元エンティティ
//
// return:
//   - dto.SenderForm: 送信者フォーム DTO
func SenderFromEntity(e entity.Sender) dto.SenderForm {
	f := dto.SenderForm{
		ID:             e.ID,
		Name:           e.Name,
		DepartmentName: e.DepartmentName,
	}
	for _, st := range e.SenderTalks {
		f.SenderTalks = append(f.SenderTalks, SenderTalkFromEntityNoSender(st))
	}
	return f
}

// senderFromEntityShallow は entity.Sender を SenderTalks なしの dto.SenderForm に変換する。
//
// args:
//   - e entity.Sender: 変換元エンティティ
//
// return:
//   - dto.SenderForm: 送信者フォーム DTO（SenderTalks なし）
func senderFromEntityShallow(e entity.Sender) dto.SenderForm {
	return dto.SenderForm{
		ID:             e.ID,
		Name:           e.Name,
		DepartmentName: e.DepartmentName,
	}
}

// SenderToEntity は dto.SenderForm を entity.Sender に変換する。
//
// args:
//   - f dto.SenderForm: 変換元フォーム DTO
//
// return:
//   - entity.Sender: 送信者エンティティ
func SenderToEntity(f dto.SenderForm) entity.Sender {
	e := entity.Sender{
		ID:             f.ID,
		Name:           f.Name,
		DepartmentName: f.DepartmentName,
	}
	for _, stf := range f.SenderTalks {
		e.SenderTalks = append(e.SenderTalks, SenderTalkToEntity(stf))
	}
	return e
}

// SenderTalkFromEntity は entity.SenderTalk を dto.SenderTalkForm に変換する。Sender があればネストする。
//
// args:
//   - e entity.SenderTalk: 変換元エンティティ
//
// return:
//   - dto.SenderTalkForm: 送信者トークフォーム DTO
func SenderTalkFromEntity(e entity.SenderTalk) dto.SenderTalkForm {
	f := SenderTalkFromEntityNoSender(e)
	if e.Sender.ID != 0 {
		s := senderFromEntityShallow(e.Sender)
		f.Sender = &s
	}
	return f
}

// SenderTalkFromEntityNoSender は entity.SenderTalk を Sender なしの dto.SenderTalkForm に変換する。
//
// args:
//   - e entity.SenderTalk: 変換元エンティティ
//
// return:
//   - dto.SenderTalkForm: 送信者トークフォーム DTO（Sender なし）
func SenderTalkFromEntityNoSender(e entity.SenderTalk) dto.SenderTalkForm {
	return dto.SenderTalkForm{
		ID:         e.ID,
		Content:    e.Content,
		SenderID:   strconv.FormatInt(e.SenderID, 10),
		QuestionID: strconv.FormatInt(e.QuestionID, 10),
		CreatedAt:  timeToISO(e.CreatedAt),
		UpdatedAt:  timeToISO(e.UpdatedAt),
		DeletedAt:  deletedAtToISO(e.DeletedAt),
	}
}

// SenderTalkToEntity は dto.SenderTalkForm を entity.SenderTalk に変換する。
//
// args:
//   - f dto.SenderTalkForm: 変換元フォーム DTO
//
// return:
//   - entity.SenderTalk: 送信者トークエンティティ
func SenderTalkToEntity(f dto.SenderTalkForm) entity.SenderTalk {
	e := entity.SenderTalk{
		ID:         f.ID,
		Content:    f.Content,
		SenderID:   f.SenderIDInt64(),
		QuestionID: f.QuestionIDInt64(),
	}
	if f.Sender != nil {
		e.Sender = SenderToEntity(*f.Sender)
		if e.SenderID <= 0 {
			e.SenderID = e.Sender.ID
		}
	}
	return e
}
