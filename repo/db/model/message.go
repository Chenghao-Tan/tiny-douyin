package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content    string
	FromUserID uint
	ToUserID   uint
	ToUser     *User `gorm:"foreignKey:ToUserID"`
}
