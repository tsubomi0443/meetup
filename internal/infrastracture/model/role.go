package model

// Role maps to table roles.
type Role struct {
	ID       uint64 `gorm:"primaryKey"`
	RoleName string `gorm:"size:50;not null"`

	Users []User `gorm:"foreignKey:RoleID"`
}
