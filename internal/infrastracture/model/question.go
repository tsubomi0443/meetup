package model

import "time"

// Question maps to table questions.
type Question struct {
	ID               uint64 `gorm:"primaryKey"`
	OriginQuestionID *uint64 `gorm:"index:idx_questions_origin_question_id"`
	SupportID        *uint64 `gorm:"index:idx_questions_support_id"`
	Title            string  `gorm:"size:255"`
	Content          string  `gorm:"type:text;not null"`
	Due              *time.Time
	CreatedAt        time.Time `gorm:"not null"`

	OriginQuestion *Question  `gorm:"foreignKey:OriginQuestionID"`
	SubQuestions   []Question `gorm:"foreignKey:OriginQuestionID"`
	Support        *Support   `gorm:"foreignKey:SupportID"`

	Answers []Answer `gorm:"foreignKey:QuestionID"`
	Memos   []Memo   `gorm:"foreignKey:QuestionID"`
	Tags    []Tag    `gorm:"many2many:tag_managers"`

	EscalationsFrom []Escalation `gorm:"foreignKey:FromQuestionID"`
	EscalationsTo   []Escalation `gorm:"foreignKey:ToQuestionID"`
}
