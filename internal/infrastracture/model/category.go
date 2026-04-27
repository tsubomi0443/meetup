package model

// Category maps to table categories.
type Category struct {
	ID           uint64 `gorm:"primaryKey"`
	CategoryName string `gorm:"size:255;not null"`

	Tags []Tag `gorm:"foreignKey:CategoryID"`
}
