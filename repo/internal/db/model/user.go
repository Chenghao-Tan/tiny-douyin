package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" redis:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;precision:0;index" redis:"createdat"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;precision:0" redis:"updatedat"`
	DeletedAt gorm.DeletedAt `gorm:"index" redis:"-"`

	Username       string    `gorm:"size:32;uniqueIndex" redis:"username"`
	Password       string    `gorm:"size:64" redis:"-"` // bcrypt结果长度为60
	Signature      string    `gorm:"size:256" redis:"signature"`
	Works          []Video   `gorm:"foreignKey:AuthorID" redis:"-"`
	WorksCount     uint      `gorm:"default:0" redis:"-"`
	Favorites      []*Video  `gorm:"many2many:favorite" redis:"-"`
	FavoritesCount uint      `gorm:"default:0" redis:"-"`
	FavoritedCount uint      `gorm:"default:0" redis:"-"`
	Comments       []Comment `gorm:"foreignKey:AuthorID" redis:"-"`
	CommentsCount  uint      `gorm:"default:0" redis:"-"`
	Follows        []*User   `gorm:"many2many:follow;joinForeignKey:user_id;joinReferences:follow_id" redis:"-"`
	FollowsCount   uint      `gorm:"default:0" redis:"-"`
	Followers      []*User   `gorm:"many2many:follow;joinForeignKey:follow_id;joinReferences:user_id" redis:"-"`
	FollowersCount uint      `gorm:"default:0" redis:"-"`
	Messages       []Message `gorm:"foreignKey:FromUserID" redis:"-"`
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
