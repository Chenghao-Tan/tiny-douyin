package db

import (
	"douyin/repo/internal/db/model"

	"context"
	"time"

	"gorm.io/gorm"
)

// 获取视频主键最大值
func MaxVideoID(ctx context.Context) (max uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Video{}).Select("IFNULL(MAX(id),0)").Scan(&max).Error
	return max, err
}

// 创建视频
func CreateVideo(ctx context.Context, authorID uint, title string) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		video = &model.Video{Title: title, AuthorID: authorID}
		author := &model.User{ID: authorID}

		err2 := tx.Model(&model.Video{}).Create(video).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(author).Update("WorksCount", gorm.Expr("works_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return video, nil
}

// 删除视频
func DeleteVideo(ctx context.Context, id uint, permanently bool) (err error) {
	DB := _db.WithContext(ctx)
	return DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		video := &model.Video{ID: id}

		var results []model.Video
		err2 := tx.Model(&model.Video{}).Select("id").Where("id=?", id).Limit(1).Find(&results).Error
		if err2 != nil {
			return err2
		}
		if len(results) == 0 { // 不允许凭空删除
			return ErrorRecordNotExists
		}

		var authorID uint
		err2 = tx.Model(video).Select("author_id").Scan(&authorID).Error
		if err2 != nil {
			return err2
		}
		author := &model.User{ID: authorID}

		if permanently {
			err2 = tx.Model(&model.Video{}).Unscoped().Delete(video).Error
		} else {
			err2 = tx.Model(&model.Video{}).Delete(video).Error
		}
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(author).Update("WorksCount", gorm.Expr("works_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
}

// 根据创建时间查找视频列表(num==-1时取消数量限制) (select: ID)
func FindVideosByCreatedAt(ctx context.Context, createdAt int64, forward bool, num int) (videoIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.Video{}).Select("id").Where("created_at>?", stop).Order("created_at").Limit(num).Find(&videoIDs).Error
	} else {
		err = DB.Model(&model.Video{}).Select("id").Where("created_at<?", stop).Order("created_at desc").Limit(num).Find(&videoIDs).Error
	}
	if err != nil {
		return []uint{}, err
	}
	return videoIDs, nil
}

// 读取视频基本信息 (select: ID, CreatedAt, UpdatedAt, Title, AuthorID)
func ReadVideoBasics(ctx context.Context, id uint) (video *model.Video, err error) {
	DB := _db.WithContext(ctx)
	video = &model.Video{}
	err = DB.Model(&model.Video{}).Select("id", "created_at", "updated_at", "title", "author_id").Where("id=?", id).First(video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

// 读取点赞(用户)列表 (select: Favorited.ID)
func ReadVideoFavorited(ctx context.Context, id uint) (userIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Video{ID: id}).Select("id").Association("Favorited").Find(&userIDs)
	if err != nil {
		return []uint{}, err
	}
	return userIDs, nil
}

// 读取点赞(用户)数量
func CountVideoFavorited(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.Video{ID: id}).Select("FavoritedCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 读取评论列表(num==-1时取消数量限制) (select: Comments.ID)
func ReadVideoComments(ctx context.Context, id uint, createdAt int64, forward bool, num int) (commentIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.Video{ID: id}).Select("id").Where("created_at>?", stop).Order("created_at").Limit(num).Association("Comments").Find(&commentIDs)
	} else {
		err = DB.Model(&model.Video{ID: id}).Select("id").Where("created_at<?", stop).Order("created_at desc").Limit(num).Association("Comments").Find(&commentIDs)
	}
	if err != nil {
		return []uint{}, err
	}
	return commentIDs, nil
}

// 读取评论数量
func CountVideoComments(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.Video{ID: id}).Select("CommentsCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 检查评论所属
func CheckVideoComments(ctx context.Context, id uint, commentID uint) (isIts bool) {
	DB := _db.WithContext(ctx)
	var results []model.Comment
	err := DB.Model(&model.Video{ID: id}).Select("id").Where("id=?", commentID).Limit(1).Association("Comments").Find(&results)
	return err == nil && len(results) > 0
}
