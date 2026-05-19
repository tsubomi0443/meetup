package dto

import "strconv"

// =====================
// QuestionForm は質問のフォーム DTO である。DeletedAt は論理削除日時を表す。
type QuestionForm struct {
	ID               int64                 `json:"id"`
	OriginQuestionID *string               `json:"originQuestionId,omitempty"`
	SupportID        *int64                `json:"supportId,omitempty"`
	Title            string                `json:"title"`
	Content          string                `json:"content"`
	Due              *string               `json:"due,omitempty"`
	CreatedAt        *string               `json:"createdAt,omitempty"`
	UpdatedAt        *string               `json:"updatedAt,omitempty"`
	DeletedAt        *string               `json:"deletedAt,omitempty"`
	OriginQuestion   *QuestionForm         `json:"originQuestion,omitempty"`
	SubQuestions     []QuestionForm        `json:"subQuestions,omitempty"`
	Support          *SupportForm          `json:"support,omitempty"`
	Answers          []AnswerForm          `json:"answers,omitempty"`
	Memos            []MemoForm            `json:"memos,omitempty"`
	Tags             []TagForm             `json:"tags,omitempty"`
	EscalationsFrom  []EscalationForm      `json:"escalationsFrom,omitempty"`
	EscalationsTo    []EscalationForm      `json:"escalationsTo,omitempty"`
	RelatedQuestions []RelatedQuestionForm `json:"relatedQuestions,omitempty"`
	SenderTalks      []SenderTalkForm      `json:"senderTalks,omitempty"`
}

// OriginQuestionIDInt64 は OriginQuestionID 文字列を int64 に変換する。nil または空文字の場合は -1 を返す。
//
// return:
//   - int64: 元質問 ID（未設定・変換失敗時は -1）
func (f QuestionForm) OriginQuestionIDInt64() int64 {
	if f.OriginQuestionID == nil || *f.OriginQuestionID == "" {
		return -1
	}
	if val, err := strconv.ParseInt(*f.OriginQuestionID, 10, 64); err == nil {
		return val
	}
	return -1
}
