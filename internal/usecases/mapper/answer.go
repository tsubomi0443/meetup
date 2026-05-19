package mapper

import (
	"strconv"
	"strings"
	"time"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// 回答（entity.Answer）

// =====================
// AnswerFromEntity は entity.Answer を dto.AnswerForm に変換する。
//
// args:
//   - e entity.Answer: 変換元エンティティ
//
// return:
//   - dto.AnswerForm: 回答フォーム DTO
func AnswerFromEntity(e entity.Answer) dto.AnswerForm {
	f := dto.AnswerForm{
		ID:        e.ID,
		UserID:    strconv.FormatInt(e.UserID, 10),
		Content:   e.Content,
		IsFinal:   e.IsFinal,
		CreatedAt: timeToISO(e.CreatedAt),
		UpdatedAt: timeToISO(e.UpdatedAt),
		DeletedAt: deletedAtToISO(e.DeletedAt),
	}
	if e.User.ID != 0 {
		u := UserFromEntity(e.User)
		f.User = &u
	}
	for _, rm := range e.ReferManagers {
		if rm.Refer.ID != 0 {
			f.Refers = append(f.Refers, referFromEntityShallow(rm.Refer))
		}
	}
	return f
}

// AnswerToEntity は dto.AnswerForm を entity.Answer に変換する。CreatedAt 未指定時は現在時刻を設定する。
//
// args:
//   - f dto.AnswerForm: 変換元フォーム DTO
//
// return:
//   - entity.Answer: 回答エンティティ
func AnswerToEntity(f dto.AnswerForm) entity.Answer {
	e := entity.Answer{
		ID:      f.ID,
		UserID:  f.UserIDInt64(),
		Content: f.Content,
		IsFinal: f.IsFinal,
	}
	if f.CreatedAt == nil || (f.CreatedAt != nil && *f.CreatedAt == "") {
		e.CreatedAt = time.Now()
	} else {
		e.CreatedAt = isoToTime(f.CreatedAt)
	}
	if f.User != nil {
		e.User = UserToEntity(*f.User)
	}
	answerID := f.ID
	for _, rf := range f.Refers {
		rm := entity.ReferManager{
			AnswerID: answerID,
			ReferID:  rf.ID,
		}
		if rf.ID != 0 {
			rm.Refer = entity.Refer{ID: rf.ID}
		} else if strings.TrimSpace(rf.Title) != "" && strings.TrimSpace(rf.URL) != "" {
			rm.Refer = ReferToEntity(rf)
		}
		e.ReferManagers = append(e.ReferManagers, rm)
	}
	return e
}
