package mapper

import (
	"strconv"
	"time"

	"meetup/internal/domains/entity"
	"meetup/internal/usecases/dto"
)

// =====================
// エスカレーション（entity.Escalation）

// =====================
// EscalationFromEntity は entity.Escalation を dto.EscalationForm に変換する。
//
// args:
//   - e entity.Escalation: 変換元エンティティ
//
// return:
//   - dto.EscalationForm: エスカレーションフォーム DTO
func EscalationFromEntity(e entity.Escalation) dto.EscalationForm {
	return dto.EscalationForm{
		ID:             e.ID,
		FromQuestionID: strconv.FormatInt(e.FromQuestionID, 10),
		ToQuestionID:   strconv.FormatInt(e.ToQuestionID, 10),
		EscalatedAt:    timeToISO(e.EscalatedAt),
		CreatedAt:      timeToISO(e.CreatedAt),
		UpdatedAt:      timeToISO(e.UpdatedAt),
		DeletedAt:      deletedAtToISO(e.DeletedAt),
	}
}

// EscalationToEntity は dto.EscalationForm を entity.Escalation に変換する。EscalatedAt 未指定時は現在時刻を設定する。
//
// args:
//   - f dto.EscalationForm: 変換元フォーム DTO
//
// return:
//   - entity.Escalation: エスカレーションエンティティ
func EscalationToEntity(f dto.EscalationForm) entity.Escalation {
	e := entity.Escalation{
		ID:             f.ID,
		FromQuestionID: f.FromQuestionIDInt64(),
		ToQuestionID:   f.ToQuestionIDInt64(),
	}
	if f.EscalatedAt == nil || *f.EscalatedAt == "" {
		e.EscalatedAt = time.Now()
	} else {
		e.EscalatedAt = isoToTime(f.EscalatedAt)
	}
	return e
}
