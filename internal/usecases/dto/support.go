package dto

import "strconv"

// =====================
// SupportForm はサポート（対応）のフォーム DTO である。
type SupportForm struct {
	ID              int64              `json:"id"`
	UserID          string             `json:"userId"`
	SupportStatusID string             `json:"supportStatusId"`
	CreatedAt       *string            `json:"createdAt,omitempty"`
	UpdatedAt       *string            `json:"updatedAt,omitempty"`
	DeletedAt       *string            `json:"deletedAt,omitempty"`
	User            *UserForm          `json:"user,omitempty"`
	SupportStatus   *SupportStatusForm `json:"supportStatus,omitempty"`
	Question        *QuestionForm      `json:"question,omitempty"`
}

// UserIDInt64 は UserID 文字列を int64 に変換する。
//
// return:
//   - int64: ユーザー ID（変換失敗時は -1）
func (f SupportForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}

// SupportStatusIDInt64 は SupportStatusID 文字列を int64 に変換する。
//
// return:
//   - int64: サポートステータス ID（変換失敗時は -1）
func (f SupportForm) SupportStatusIDInt64() int64 {
	if val, err := strconv.ParseInt(f.SupportStatusID, 10, 64); err == nil {
		return val
	}
	return -1
}
