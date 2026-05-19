package dto

import "strconv"

// =====================
// EscalationForm は質問エスカレーションのフォーム DTO である。
type EscalationForm struct {
	ID             int64         `json:"id"`
	FromQuestionID string        `json:"fromQuestionId"`
	ToQuestionID   string        `json:"toQuestionId"`
	EscalatedAt    *string       `json:"escalatedAt,omitempty"`
	CreatedAt      *string       `json:"createdAt,omitempty"`
	UpdatedAt      *string       `json:"updatedAt,omitempty"`
	DeletedAt      *string       `json:"deletedAt,omitempty"`
	FromQuestion   *QuestionForm `json:"fromQuestion,omitempty"`
	ToQuestion     *QuestionForm `json:"toQuestion,omitempty"`
}

// FromQuestionIDInt64 は FromQuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: エスカレーション元質問 ID（変換失敗時は -1）
func (f EscalationForm) FromQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.FromQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// ToQuestionIDInt64 は ToQuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: エスカレーション先質問 ID（変換失敗時は -1）
func (f EscalationForm) ToQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.ToQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}
