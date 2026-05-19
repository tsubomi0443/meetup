package dto

import "strconv"

// =====================
// RelatedQuestionForm は関連質問の紐付けフォーム DTO である。
type RelatedQuestionForm struct {
	ID                int64         `json:"id"`
	QuestionID        string        `json:"questionId"`
	RelatedQuestionID string        `json:"relatedQuestionId"`
	CreatedAt         *string       `json:"createdAt,omitempty"`
	UpdatedAt         *string       `json:"updatedAt,omitempty"`
	DeletedAt         *string       `json:"deletedAt,omitempty"`
	Question          *QuestionForm `json:"question,omitempty"`
	RelatedQuestion   *QuestionForm `json:"relatedQuestion,omitempty"`
}

// QuestionIDInt64 は QuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 質問 ID（変換失敗時は -1）
func (f RelatedQuestionForm) QuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.QuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}

// RelatedQuestionIDInt64 は RelatedQuestionID 文字列を int64 に変換する。
//
// return:
//   - int64: 関連質問 ID（変換失敗時は -1）
func (f RelatedQuestionForm) RelatedQuestionIDInt64() int64 {
	if val, err := strconv.ParseInt(f.RelatedQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}
