package dto

import "strconv"

// =====================
// ReferManagerForm は回答と参照リンクの紐付けフォーム DTO である。
type ReferManagerForm struct {
	ID        int64       `json:"id"`
	AnswerID  string      `json:"answerId"`
	ReferID   string      `json:"referId"`
	CreatedAt *string     `json:"createdAt,omitempty"`
	UpdatedAt *string     `json:"updatedAt,omitempty"`
	DeletedAt *string     `json:"deletedAt,omitempty"`
	Answer    *AnswerForm `json:"answer,omitempty"`
	Refer     *ReferForm  `json:"refer,omitempty"`
}

// AnswerIDInt64 は AnswerID 文字列を int64 に変換する。
//
// return:
//   - int64: 回答 ID（変換失敗時は -1）
func (f ReferManagerForm) AnswerIDInt64() int64 {
	if val, err := strconv.ParseInt(f.AnswerID, 10, 64); err == nil {
		return val
	}
	return -1
}

// ReferIDInt64 は ReferID 文字列を int64 に変換する。
//
// return:
//   - int64: 参照リンク ID（変換失敗時は -1）
func (f ReferManagerForm) ReferIDInt64() int64 {
	if val, err := strconv.ParseInt(f.ReferID, 10, 64); err == nil {
		return val
	}
	return -1
}
