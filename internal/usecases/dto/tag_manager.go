package dto

import "strconv"

// =====================
// TagManagerForm は質問とタグの紐付けフォーム DTO である。
type TagManagerForm struct {
	ID         int64         `json:"id"`
	TagID      string        `json:"tagId"`
	QuestionID string        `json:"questionId"`
	CreatedAt  *string       `json:"createdAt,omitempty"`
	UpdatedAt  *string       `json:"updatedAt,omitempty"`
	DeletedAt  *string       `json:"deletedAt,omitempty"`
	Tag        *TagForm      `json:"tag,omitempty"`
	Question   *QuestionForm `json:"question,omitempty"`
}

// TagIDInt64 は TagID 文字列を int64 に変換する。
//
// return:
//   - int64: タグ ID（変換失敗時は -1）
func (f TagManagerForm) TagIDInt64() int64 {
	if val, err := strconv.ParseInt(f.TagID, 10, 64); err == nil {
		return val
	}
	return -1
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f TagManagerForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}
