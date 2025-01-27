package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" redis:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;precision:0;index" redis:"createdat"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;precision:0" redis:"updatedat"`
	DeletedAt gorm.DeletedAt `gorm:"index" redis:"-"`

	Content  string `gorm:"size:256" redis:"content"`
	AuthorID uint   `redis:"authorid"`
	VideoID  uint   `redis:"videoid"`
}
