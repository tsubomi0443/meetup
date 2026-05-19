package dto

// =====================
// ReferForm は参照リンクのフォーム DTO である。
type ReferForm struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	URL       string       `json:"url"`
	CreatedAt *string      `json:"createdAt,omitempty"`
	UpdatedAt *string      `json:"updatedAt,omitempty"`
	DeletedAt *string      `json:"deletedAt,omitempty"`
	Answers   []AnswerForm `json:"answers,omitempty"`
}
