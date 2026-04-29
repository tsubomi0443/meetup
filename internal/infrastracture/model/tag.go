package model

// Tag maps to table tags.
type Tag struct {
	ID         uint64 `gorm:"primaryKey"`
	Title      string `gorm:"size:255;not null"`
	Usage      int    `gorm:"not null;default:0"`
	CategoryID uint64 `gorm:"not null;index:idx_tags_category_id"`

	Category Category `gorm:"foreignKey:CategoryID"`

	Questions []Question `gorm:"many2many:tag_managers"`
}
