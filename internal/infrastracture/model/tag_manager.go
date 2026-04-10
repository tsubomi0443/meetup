package model

// TagManager maps to join table tag_managers (many-to-many between tags and questions).
type TagManager struct {
	ID         uint64 `gorm:"primaryKey"`
	TagID      uint64 `gorm:"not null;index:idx_tag_managers_tag_id"`
	QuestionID uint64 `gorm:"not null;index:idx_tag_managers_question_id"`

	Tag      Tag      `gorm:"foreignKey:TagID"`
	Question Question `gorm:"foreignKey:QuestionID"`
}
