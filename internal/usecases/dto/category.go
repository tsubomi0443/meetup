package dto

// =====================
// CategoryForm はカテゴリのフォーム DTO である。
type CategoryForm struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt *string   `json:"createdAt,omitempty"`
	UpdatedAt *string   `json:"updatedAt,omitempty"`
	DeletedAt *string   `json:"deletedAt,omitempty"`
	Tags      []TagForm `json:"tags,omitempty"`
}
