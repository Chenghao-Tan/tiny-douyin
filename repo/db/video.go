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

// 删除视频
func DeleteVideo(ctx context.Context, id uint, permanently bool) (err error) {
	DB := _db.WithContext(ctx)
	video := &model.Video{Model: gorm.Model{ID: id}}
	if permanently {
		err = DB.Model(&model.Video{}).Unscoped().Delete(video).Error
	} else {
		err = DB.Model(&model.Video{}).Delete(video).Error
	}
	return err
}

// 根据创建时间查找视频列表(num==-1时取消数量限制) (select: ID, CreatedAt)
func FindVideosByCreatedAt(ctx context.Context, createdAt int64, forward bool, num int) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.Video{}).Select("id", "created_at").Where("created_at>?", stop).Order("created_at").Limit(num).Find(&videos).Error
	} else {
		err = DB.Model(&model.Video{}).Select("id", "created_at").Where("created_at<?", stop).Order("created_at desc").Limit(num).Find(&videos).Error
	}
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// 读取视频基本信息 (select: *)
func ReadVideoBasics(ctx context.Context, id uint) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	video = &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", id).First(video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

// 读取点赞(用户)列表 (select: Favorited.ID)
func ReadVideoFavorited(ctx context.Context, id uint) (users []model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Select("id").Association("Favorited").Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

// 读取点赞(用户)数量
func CountVideoFavorited(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Select("FavoritedCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
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

// 计算评论数量
func CountVideoComments(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.Video{Model: gorm.Model{ID: id}}).Association("Comments").Count()
}
