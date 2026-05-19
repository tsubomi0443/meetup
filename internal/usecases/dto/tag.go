package dto

import "strconv"

// =====================
// TagForm はタグのフォーム DTO である。
type TagForm struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	Usage      int            `json:"usage"`
	CategoryID string         `json:"categoryId"`
	CreatedAt  *string        `json:"createdAt,omitempty"`
	UpdatedAt  *string        `json:"updatedAt,omitempty"`
	DeletedAt  *string        `json:"deletedAt,omitempty"`
	Category   *CategoryForm  `json:"category,omitempty"`
	Questions  []QuestionForm `json:"questions,omitempty"`
}

// CategoryIDInt64 は CategoryID 文字列を int64 に変換する。
//
// return:
//   - int64: カテゴリ ID（変換失敗時は -1）
func (tf TagForm) CategoryIDInt64() int64 {
	if val, err := strconv.ParseInt(tf.CategoryID, 10, 64); err == nil {
		return val
	}
	return -1
}
