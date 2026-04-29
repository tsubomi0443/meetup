package model

import "time"

// Escalation maps to table escalations.
type Escalation struct {
	ID             uint64    `gorm:"primaryKey"`
	FromQuestionID uint64    `gorm:"not null;index:idx_escalations_from_question_id"`
	ToQuestionID   uint64    `gorm:"not null;index:idx_escalations_to_question_id"`
	EscalatedAt    time.Time `gorm:"not null"`

	FromQuestion Question `gorm:"foreignKey:FromQuestionID"`
	ToQuestion   Question `gorm:"foreignKey:ToQuestionID"`
}
