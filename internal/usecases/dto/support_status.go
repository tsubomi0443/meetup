package dto

// =====================
// SupportStatusForm はサポートステータスのフォーム DTO である。
type SupportStatusForm struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	CreatedAt *string       `json:"createdAt,omitempty"`
	UpdatedAt *string       `json:"updatedAt,omitempty"`
	DeletedAt *string       `json:"deletedAt,omitempty"`
	Supports  []SupportForm `json:"supports,omitempty"`
}
