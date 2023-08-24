package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content  string `gorm:"size:256"`
	AuthorID uint
	VideoID  uint
}
