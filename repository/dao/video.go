package dao

import (
	"douyin/repository/model"

	"context"
	"time"

	"gorm.io/gorm/clause"
)

// 根据更新时间获取视频列表
func FindVideosByUpdatedAt(ctx context.Context, updatedAt int64, forward bool, num int) (videos []model.Video, err error) {
	DB := GetDB(ctx)
	stop := time.Unix(updatedAt, 0)
	if forward {
		err = DB.Model(&model.Video{}).Where("updated_at>?", stop).Order("updated_at").Limit(num).Preload(clause.Associations).Find(&videos).Error
	} else {
		err = DB.Model(&model.Video{}).Where("updated_at<?", stop).Order("updated_at desc").Limit(num).Preload(clause.Associations).Find(&videos).Error
	}
	return videos, err
}

// 创建Video
func CreateVideo(ctx context.Context, videoInfo *model.Video) (video *model.Video, err error) {
	DB := GetDB(ctx)
	err = DB.Model(&model.Video{}).Create(videoInfo).Error
	return videoInfo, err
}
