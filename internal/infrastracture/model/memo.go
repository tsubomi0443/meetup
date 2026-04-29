package model

// Memo maps to table memos.
type Memo struct {
	ID         uint64 `gorm:"primaryKey"`
	QuestionID uint64 `gorm:"not null;index:idx_memos_question_id"`
	UserID     uint64 `gorm:"not null;index:idx_memos_user_id"`
	Content    string `gorm:"type:text;not null"`

	Question Question `gorm:"foreignKey:QuestionID"`
	User     User     `gorm:"foreignKey:UserID"`
}
