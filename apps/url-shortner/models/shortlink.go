package models

import (
	"gorm.io/gorm"
)

type ShortLink struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Author   string `gorm:"not null"`
	URL      string `gorm:"not null"`
	ShortURL string `gorm:"uniqueIndex;not null"`
}
