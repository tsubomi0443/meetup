package model

import "time"

// Answer maps to table answers.
type Answer struct {
	ID         uint64     `gorm:"primaryKey"`
	UserID     uint64     `gorm:"not null;index:idx_answers_user_id"`
	QuestionID uint64     `gorm:"not null;index:idx_answers_question_id"`
	Content    string     `gorm:"type:text;not null"`
	AnsweredAt *time.Time
	CreatedAt  time.Time  `gorm:"not null"`

	User     User     `gorm:"foreignKey:UserID"`
	Question Question `gorm:"foreignKey:QuestionID"`

	Refers []Refer `gorm:"many2many:refer_managers"`
}
