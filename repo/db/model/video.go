package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title     string
	UserID    uint
	Favorited []*User   `gorm:"many2many:favorite"`
	Comments  []Comment `gorm:"foreignKey:VideoID"`
}
