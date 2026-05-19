package dto

import "strconv"

// =====================
// MemoForm は質問メモのフォーム DTO である。
type MemoForm struct {
	ID         int64         `json:"id"`
	QuestionID string        `json:"questionId"`
	UserID     string        `json:"userId"`
	Content    string        `json:"content"`
	CreatedAt  *string       `json:"createdAt,omitempty"`
	UpdatedAt  *string       `json:"updatedAt,omitempty"`
	DeletedAt  *string       `json:"deletedAt,omitempty"`
	Question   *QuestionForm `json:"question,omitempty"`
	User       *UserForm     `json:"user,omitempty"`
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f MemoForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// UserIDInt64 は UserID 文字列を int64 に変換する。
//
// return:
//   - int64: ユーザー ID（変換失敗時は -1）
func (f MemoForm) UserIDInt64() int64 {
	if val, err := strconv.ParseInt(f.UserID, 10, 64); err == nil {
		return val
	}
	return -1
}
