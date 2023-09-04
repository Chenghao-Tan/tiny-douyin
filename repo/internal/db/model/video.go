package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID        uint           `gorm:"primaryKey" redis:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;precision:0;index" redis:"createdat"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;precision:0" redis:"updatedat"`
	DeletedAt gorm.DeletedAt `gorm:"index" redis:"-"`

	Title          string    `gorm:"size:256" redis:"title"`
	AuthorID       uint      `redis:"authorid"`
	Favorited      []*User   `gorm:"many2many:favorite" redis:"-"`
	FavoritedCount uint      `gorm:"default:0" redis:"-"`
	Comments       []Comment `gorm:"foreignKey:VideoID" redis:"-"`
	CommentsCount  uint      `gorm:"default:0" redis:"-"`
}
