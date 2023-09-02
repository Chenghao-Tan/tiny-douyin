package model

type Message struct {
	Model
	Content    string `gorm:"size:256" redis:"content"`
	FromUserID uint   `redis:"fromuserid"`
	ToUserID   uint   `redis:"touserid"`
	ToUser     *User  `gorm:"foreignKey:ToUserID" redis:"-"`
}
