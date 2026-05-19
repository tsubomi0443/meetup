package dto

import "strconv"

// =====================
// AnswerForm は回答のフォーム DTO である。
type AnswerForm struct {
	ID        int64       `json:"id"`
	UserID    string      `json:"userId"`
	Content   string      `json:"content"`
	IsFinal   bool        `json:"isFinal"`
	CreatedAt *string     `json:"createdAt,omitempty"`
	UpdatedAt *string     `json:"updatedAt,omitempty"`
	DeletedAt *string     `json:"deletedAt,omitempty"`
	User      *UserForm   `json:"user,omitempty"`
	Refers    []ReferForm `json:"refers,omitempty"`
}

// UserIDInt64 は UserID 文字列を int64 に変換する。
//
// return:
//   - int64: ユーザー ID（変換失敗時は -1）
func (f AnswerForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}
