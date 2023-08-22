package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content  string
	AuthorID uint
	VideoID  uint
}
