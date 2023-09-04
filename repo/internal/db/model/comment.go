package model

type Comment struct {
	Model
	Content  string `gorm:"size:256" redis:"content"`
	AuthorID uint   `redis:"authorid"`
	VideoID  uint   `redis:"videoid"`
}
