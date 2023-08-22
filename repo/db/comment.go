package db

import (
	"douyin/repo/db/model"

	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 评论
func CreateComment(ctx context.Context, user_id uint, video_id uint, content string) (comment *model.Comment, err error) {
	DB := _db.WithContext(ctx)
	comment = &model.Comment{Content: content, AuthorID: user_id, VideoID: video_id}
	err = DB.Model(&model.Comment{}).Create(comment).Error
	return comment, err
}

// 删除评论
func DeleteComment(ctx context.Context, comment_id uint, permanently bool) (err error) {
	DB := _db.WithContext(ctx)
	comment := &model.Comment{Model: gorm.Model{ID: comment_id}}
	if permanently {
		err = DB.Model(&model.Comment{}).Unscoped().Delete(comment).Error
	} else {
		err = DB.Model(&model.Comment{}).Delete(comment).Error
	}
	return err
}

// 根据视频ID和发布时间获取评论列表
func FindCommentsByCreatedAt(ctx context.Context, video_id uint, forward bool) (comments []model.Comment, err error) {
	DB := _db.WithContext(ctx)
	if forward {
		err = DB.Model(&model.Comment{}).Where("video_id=?", video_id).Order("created_at").Preload(clause.Associations).Find(&comments).Error
	} else {
		err = DB.Model(&model.Comment{}).Where("video_id=?", video_id).Order("created_at desc").Preload(clause.Associations).Find(&comments).Error
	}
	return comments, err
}
