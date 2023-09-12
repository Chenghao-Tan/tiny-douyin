package db

import (
	"douyin/repo/internal/db/model"

	"context"

	"gorm.io/gorm"
)

// 获取评论主键最大值
func MaxCommentID(ctx context.Context) (max uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Comment{}).Select("IFNULL(MAX(id),0)").Find(&max).Error
	return max, err
}

// 创建评论
func CreateComment(ctx context.Context, authorID uint, videoID uint, content string) (comment *model.Comment, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		comment = &model.Comment{Content: content, AuthorID: authorID, VideoID: videoID}
		author := &model.User{ID: authorID}
		video := &model.Video{ID: videoID}

		err2 := tx.Model(&model.Comment{}).Create(comment).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(author).Update("CommentsCount", gorm.Expr("comments_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(video).Update("CommentsCount", gorm.Expr("comments_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// 删除评论
func DeleteComment(ctx context.Context, id uint, permanently bool) (err error) {
	DB := _db.WithContext(ctx)
	return DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		comment := &model.Comment{ID: id}

		var results []model.Comment
		err2 := tx.Model(&model.Comment{}).Select("id").Where("id=?", id).Limit(1).Find(&results).Error
		if err2 != nil {
			return err2
		}
		if len(results) == 0 { // 不允许凭空删除
			return ErrorRecordNotExists
		}

		var authorID uint
		err2 = tx.Model(comment).Select("author_id").Scan(&authorID).Error
		if err2 != nil {
			return err2
		}
		author := &model.User{ID: authorID}

		var videoID uint
		err2 = tx.Model(comment).Select("video_id").Scan(&videoID).Error
		if err2 != nil {
			return err2
		}
		video := &model.Video{ID: videoID}

		if permanently {
			err2 = tx.Model(&model.Comment{}).Unscoped().Delete(comment).Error
		} else {
			err2 = tx.Model(&model.Comment{}).Delete(comment).Error
		}
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(author).Update("CommentsCount", gorm.Expr("comments_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(video).Update("CommentsCount", gorm.Expr("comments_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
}

// 读取评论基本信息 (select: ID, CreatedAt, UpdatedAt, Content, AuthorID, VideoID)
func ReadCommentBasics(ctx context.Context, id uint) (comment *model.Comment, err error) {
	DB := _db.WithContext(ctx)
	comment = &model.Comment{}
	err = DB.Model(&model.Comment{}).Select("id", "created_at", "updated_at", "content", "author_id", "video_id").Where("id=?", id).First(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}
