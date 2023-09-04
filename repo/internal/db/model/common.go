package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" redis:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;precision:0;index" redis:"createdat"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;precision:0" redis:"updatedat"`
	DeletedAt gorm.DeletedAt `gorm:"index" redis:"-"`
}
