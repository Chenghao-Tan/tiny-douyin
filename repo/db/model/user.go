package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string    `gorm:"uniqueIndex;size:32"`
	Password       string    `gorm:"size:64"` // bcrypt结果长度为60
	Signature      string    `gorm:"size:256"`
	Works          []Video   `gorm:"foreignKey:AuthorID"`
	Favorites      []*Video  `gorm:"many2many:favorite"`
	FavoritesCount uint      `gorm:"default:0"`
	FavoritedCount uint      `gorm:"default:0"`
	Comments       []Comment `gorm:"foreignKey:AuthorID"`
	Follows        []*User   `gorm:"many2many:follow;joinForeignKey:user_id;joinReferences:follow_id"`
	Followers      []*User   `gorm:"many2many:follow;joinForeignKey:follow_id;joinReferences:user_id"`
	Messages       []Message `gorm:"foreignKey:FromUserID"`
}

const passwordCost = 12 //密码加密难度

// SetPassword 设置密码
func (user *User) SetPassword(password string) (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) (isValid bool) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
