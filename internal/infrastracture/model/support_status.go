package model

// SupportStatus maps to table support_statuses.
type SupportStatus struct {
	ID    uint64 `gorm:"primaryKey"`
	Title string `gorm:"size:255;not null"`

	Supports []Support `gorm:"foreignKey:SupportStatusID"`
}
