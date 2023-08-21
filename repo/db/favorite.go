package db

import (
	"douyin/repo/db/model"

	"context"

	"gorm.io/gorm"
)

// 点赞
func CreateFavorite(ctx context.Context, user_id uint, video_id uint) (err error) {
	DB := _db.WithContext(ctx)
	user := &model.User{Model: gorm.Model{ID: user_id}}
	video := &model.Video{Model: gorm.Model{ID: video_id}}
	return DB.Model(user).Association("Favorites").Append(video)
}

// 取消点赞
func DeleteFavorite(ctx context.Context, user_id uint, video_id uint) (err error) {
	DB := _db.WithContext(ctx)
	user := &model.User{Model: gorm.Model{ID: user_id}}
	video := &model.Video{Model: gorm.Model{ID: video_id}}
	return DB.Model(user).Association("Favorites").Delete(video)
}

// 检查是否点赞
func CheckFavorite(ctx context.Context, user_id uint, video_id uint) (is_favorite bool) {
	DB := _db.WithContext(ctx)
	user := &model.User{Model: gorm.Model{ID: user_id}}
	return DB.Model(user).Where("id=?", video_id).Association("Favorites").Count() > 0
}
