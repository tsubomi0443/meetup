package model

// Support maps to table supports.
type Support struct {
	ID              uint64 `gorm:"primaryKey"`
	UserID          uint64 `gorm:"not null;index:idx_supports_user_id"`
	SupportStatusID uint64 `gorm:"not null;index:idx_supports_support_status_id"`

	User          User          `gorm:"foreignKey:UserID"`
	SupportStatus SupportStatus `gorm:"foreignKey:SupportStatusID"`

	Questions []Question `gorm:"foreignKey:SupportID"`
}
