package dto

import "strconv"

// =====================
// UserForm はユーザーのフォーム DTO である。
type UserForm struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Memo      string        `json:"memo"`
	RoleID    string        `json:"roleId"`
	CreatedAt *string       `json:"createdAt,omitempty"`
	UpdatedAt *string       `json:"updatedAt,omitempty"`
	DeletedAt *string       `json:"deletedAt,omitempty"`
	Role      *RoleForm     `json:"role,omitempty"`
	Supports  []SupportForm `json:"supports,omitempty"`
	Answers   []AnswerForm  `json:"answers,omitempty"`
	Memos     []MemoForm    `json:"memos,omitempty"`
	Password  string        `json:"pass,omitempty"`
}

// RoleIDInt64 は RoleID 文字列を int64 に変換する。
//
// return:
//   - int64: ロール ID（変換失敗時は -1）
func (uf UserForm) RoleIDInt64() int64 {
	if val, err := strconv.ParseInt(uf.RoleID, 10, 64); err == nil {
		return val
	}
	return -1
}
