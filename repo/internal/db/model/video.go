package model

type Video struct {
	Model
	Title          string    `gorm:"size:256" redis:"title"`
	AuthorID       uint      `redis:"authorid"`
	Favorited      []*User   `gorm:"many2many:favorite" redis:"-"`
	FavoritedCount uint      `gorm:"default:0" redis:"-"`
	Comments       []Comment `gorm:"foreignKey:VideoID" redis:"-"`
	CommentsCount  uint      `gorm:"default:0" redis:"-"`
}
