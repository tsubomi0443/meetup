package model

// ReferManager maps to join table refer_managers (many-to-many between answers and refers).
type ReferManager struct {
	ID       uint64 `gorm:"primaryKey"`
	AnswerID uint64 `gorm:"not null;index:idx_refer_managers_answer_id"`
	ReferID  uint64 `gorm:"not null;index:idx_refer_managers_refer_id"`

	Answer Answer `gorm:"foreignKey:AnswerID"`
	Refer  Refer  `gorm:"foreignKey:ReferID"`
}
