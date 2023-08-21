package db

import (
	"douyin/repo/db/model"

	"context"
	"time"

	"gorm.io/gorm/clause"
)

// 创建Video
func CreateVideo(ctx context.Context, videoInfo *model.Video) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Video{}).Create(videoInfo).Error
	return videoInfo, err
}

// 根据视频ID找到视频
func FindVideoByID(ctx context.Context, id uint) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	video = &model.Video{}
	err = DB.Model(&model.Video{}).Where("id=?", id).Preload(clause.Associations).First(video).Error
	return video, err
}

// 根据更新时间获取视频列表
func FindVideosByUpdatedAt(ctx context.Context, updatedAt int64, forward bool, num int) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(updatedAt, 0)
	if forward {
		err = DB.Model(&model.Video{}).Where("updated_at>?", stop).Order("updated_at").Limit(num).Preload(clause.Associations).Find(&videos).Error
	} else {
		err = DB.Model(&model.Video{}).Where("updated_at<?", stop).Order("updated_at desc").Limit(num).Preload(clause.Associations).Find(&videos).Error
	}
	return videos, err
}
