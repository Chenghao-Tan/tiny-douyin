package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title          string `gorm:"size:256"`
	AuthorID       uint
	Favorited      []*User   `gorm:"many2many:favorite"`
	FavoritedCount uint      `gorm:"default:0"`
	Comments       []Comment `gorm:"foreignKey:VideoID"`
	CommentsCount  uint      `gorm:"default:0"`
}
