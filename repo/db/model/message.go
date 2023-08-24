package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content    string `gorm:"size:256"`
	FromUserID uint
	ToUserID   uint
	ToUser     *User `gorm:"foreignKey:ToUserID"`
}
