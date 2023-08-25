package db

import (
	"douyin/repo/db/model"

	"context"
	"time"

	"gorm.io/gorm"
)

// 创建评论
func CreateComment(ctx context.Context, authorID uint, videoID uint, content string) (comment *model.Comment, err error) {
	DB := _db.WithContext(ctx)
	comment = &model.Comment{Content: content, AuthorID: authorID, VideoID: videoID}
	err = DB.Model(&model.Comment{}).Create(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// 删除评论
func DeleteComment(ctx context.Context, id uint, permanently bool) (err error) {
	DB := _db.WithContext(ctx)
	comment := &model.Comment{Model: gorm.Model{ID: id}}
	if permanently {
		err = DB.Model(&model.Comment{}).Unscoped().Delete(comment).Error
	} else {
		err = DB.Model(&model.Comment{}).Delete(comment).Error
	}
	return err
}

// 根据视频ID和创建时间查找评论列表(num==-1时取消数量限制) (select: ID, CreatedAt)
func FindCommentsByCreatedAt(ctx context.Context, videoID uint, createdAt int64, forward bool, num int) (comments []model.Comment, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.Comment{}).Select("id", "created_at").Where("video_id=?", videoID).Where("created_at>?", stop).Order("created_at").Limit(num).Find(&comments).Error
	} else {
		err = DB.Model(&model.Comment{}).Select("id", "created_at").Where("video_id=?", videoID).Where("created_at<?", stop).Order("created_at desc").Limit(num).Find(&comments).Error
	}
	if err != nil {
		return comments, err
	}
	return comments, err
}

// 读取评论基本信息 (select: *)
func ReadCommentBasics(ctx context.Context, id uint) (comment *model.Comment, err error) {
	DB := _db.WithContext(ctx)
	comment = &model.Comment{}
	err = DB.Model(&model.Comment{}).Where("id=?", id).First(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}
