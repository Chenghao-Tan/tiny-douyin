package db

import (
	"douyin/repo/db/model"

	"context"
	"time"

	"gorm.io/gorm"
)

// 创建视频
func CreateVideo(ctx context.Context, authorID uint, title string) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	video = &model.Video{Title: title, AuthorID: authorID}
	err = DB.Model(&model.Video{}).Create(video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

// 根据视频ID查找视频 (select: *)
func FindVideoByID(ctx context.Context, id uint) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	video = &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", id).First(video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

// 根据更新时间查找视频列表 (select: ID, UpdatedAt)
func FindVideosByUpdatedAt(ctx context.Context, updatedAt int64, forward bool, num int) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(updatedAt, 0)
	if forward {
		err = DB.Model(&model.Video{}).Select("id", "updated_at").Where("updated_at>?", stop).Order("updated_at").Limit(num).Find(&videos).Error
	} else {
		err = DB.Model(&model.Video{}).Select("id", "updated_at").Where("updated_at<?", stop).Order("updated_at desc").Limit(num).Find(&videos).Error
	}
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// 读取点赞用户列表 (select: Favorited.ID)
func ReadVideoFavorited(ctx context.Context, id uint) (users []model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Select("id").Association("Favorited").Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

// 读取点赞用户数量
func CountVideoFavorited(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Association("Favorited").Count()
}

// 读取评论列表 (select: Comments.ID)
func ReadVideoComments(ctx context.Context, id uint) (comments []model.Comment, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Select("id").Association("Comments").Find(&comments)
	if err != nil {
		return comments, err
	}
	return comments, nil
}

// 读取评论数量
func CountVideoComments(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Association("Comments").Count()
}
