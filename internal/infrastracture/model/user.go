package model

// User maps to table users.
type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	Password string `gorm:"column:passwordd;size:255;not null"`
	Email    string `gorm:"size:255;not null"`
	RoleID   uint64 `gorm:"not null;index:idx_users_role_id"`

	Role Role `gorm:"foreignKey:RoleID"`

	Supports []Support `gorm:"foreignKey:UserID"`
	Answers  []Answer  `gorm:"foreignKey:UserID"`
	Memos    []Memo    `gorm:"foreignKey:UserID"`
}
