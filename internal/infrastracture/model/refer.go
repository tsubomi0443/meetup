package model

// Refer maps to table refers.
type Refer struct {
	ID    uint64 `gorm:"primaryKey"`
	Title string `gorm:"size:255;not null"`
	URL   string `gorm:"type:text;not null"`

	Answers []Answer `gorm:"many2many:refer_managers"`
}
