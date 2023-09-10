package db

import (
	"douyin/repo/internal/db/model"

	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// 自定义错误类型
var ErrorSelfFollow = errors.New("禁止自己关注自己")

// 获取用户主键最大值
func MaxUserID(ctx context.Context) (max uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{}).Select("IFNULL(MAX(id),0)").Scan(&max).Error
	return max, err
}

// 创建用户
func CreateUser(ctx context.Context, username string, password string, signature string) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	user = &model.User{Username: username, Signature: signature}
	err = user.SetPassword(password)
	if err != nil {
		return nil, err
	}
	err = DB.Model(&model.User{}).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 检查用户名是否可用
func CheckUserRegister(ctx context.Context, username string) (isAvailable bool) {
	DB := _db.WithContext(ctx)
	var results []model.User
	err := DB.Model(&model.User{}).Select("username").Where("username=?", username).Limit(1).Find(&results).Error
	return err == nil && len(results) == 0
}

// 检查用户名和密码是否有效
func CheckUserLogin(ctx context.Context, username string, password string) (id uint, isValid bool) {
	DB := _db.WithContext(ctx)
	var results []model.User
	err := DB.Model(&model.User{}).Select("id", "username", "password").Where("username=?", username).Limit(1).Find(&results).Error
	if err != nil || len(results) == 0 {
		return 0, false
	}
	if !results[0].CheckPassword(password) {
		return 0, false
	}
	return results[0].ID, true
}

// 读取用户基本信息 (select: ID, CreatedAt, UpdatedAt, Username, Signature)
func ReadUserBasics(ctx context.Context, id uint) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Select("id", "created_at", "updated_at", "username", "signature").Where("id=?", id).First(user).Error
	return user, err
}

// 读取作品(视频)列表 (select: Works.ID)
func ReadUserWorks(ctx context.Context, id uint) (videoIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Works").Find(&videoIDs)
	if err != nil {
		return []uint{}, err
	}
	return videoIDs, nil
}

// 读取作品(视频)数量
func CountUserWorks(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.User{ID: id}).Select("WorksCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 创建点赞关系
func CreateUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	DB := _db.WithContext(ctx)
	return DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		user := &model.User{ID: id}
		video := &model.Video{ID: videoID}

		var results []model.Video
		err2 := tx.Model(user).Select("id").Where("id=?", videoID).Limit(1).Association("Favorites").Find(&results)
		if err2 != nil {
			return err2
		}
		if len(results) > 0 { // 不允许重复创建
			return ErrorRecordExists
		}

		var authorID uint
		err2 = tx.Model(video).Select("author_id").Scan(&authorID).Error
		if err2 != nil {
			return err2
		}
		author := &model.User{ID: authorID}

		err2 = tx.Model(user).Association("Favorites").Append(video)
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(user).Update("FavoritesCount", gorm.Expr("favorites_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(author).Update("FavoritedCount", gorm.Expr("favorited_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(video).Update("FavoritedCount", gorm.Expr("favorited_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
}

// 创建点赞关系(批量)
func CreateUserFavoritesBatch(ctx context.Context, ids []uint, videoIDs []uint) (successCount int64) {
	successCount = 0
	if len(ids) != len(videoIDs) {
		return 0
	}

	DB := _db.WithContext(ctx)
	err := DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		users := make([]model.User, 0, len(ids))
		videos := make([]model.Video, 0, len(videoIDs))
		for i, id := range ids {
			var results []model.Video
			if tx.Model(&model.User{ID: id}).Select("id").Where("id=?", videoIDs[i]).Limit(1).Association("Favorites").Find(&results) == nil && len(results) == 0 { // 不允许重复创建
				users = append(users, model.User{ID: id})
				videos = append(videos, model.Video{ID: videoIDs[i]})
				successCount++
			}
		}

		err2 := tx.Model(users).Association("Favorites").Append(videos)
		if err2 != nil {
			return err2
		}

		var userBatch []model.User
		err2 = tx.Model(users).FindInBatches(&userBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, user := range userBatch {
				err3 := bx.Model(user).Update("FavoritesCount", gorm.Expr("favorites_count+?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		var videoBatch []model.Video
		err2 = tx.Model(videos).FindInBatches(&videoBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, video := range videoBatch {
				var authorID uint
				err3 := bx.Model(video).Select("author_id").Find(&authorID).Error
				if err3 != nil {
					return err3
				}
				author := &model.User{ID: authorID}

				err3 = bx.Model(author).Update("FavoritedCount", gorm.Expr("favorited_count+?", 1)).Error
				if err3 != nil {
					return err3
				}

				err3 = bx.Model(video).Update("FavoritedCount", gorm.Expr("favorited_count+?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return 0 // 因事务回滚
	}

	return successCount
}

// 删除点赞关系
func DeleteUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	DB := _db.WithContext(ctx)
	return DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		user := &model.User{ID: id}
		video := &model.Video{ID: videoID}

		var results []model.Video
		err2 := tx.Model(user).Select("id").Where("id=?", videoID).Limit(1).Association("Favorites").Find(&results)
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

		err2 = tx.Model(user).Association("Favorites").Delete(video)
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(user).Update("FavoritesCount", gorm.Expr("favorites_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(author).Update("FavoritedCount", gorm.Expr("favorited_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(video).Update("FavoritedCount", gorm.Expr("favorited_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
}

// 删除点赞关系(批量)
func DeleteUserFavoritesBatch(ctx context.Context, ids []uint, videoIDs []uint) (successCount int64) {
	successCount = 0
	if len(ids) != len(videoIDs) {
		return 0
	}

	DB := _db.WithContext(ctx)
	err := DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		users := make([]model.User, 0, len(ids))
		videos := make([]model.Video, 0, len(videoIDs))
		for i, id := range ids {
			var results []model.Video
			if tx.Model(&model.User{ID: id}).Select("id").Where("id=?", videoIDs[i]).Limit(1).Association("Favorites").Find(&results) == nil && len(results) > 0 { // 不允许凭空删除
				users = append(users, model.User{ID: id})
				videos = append(videos, model.Video{ID: videoIDs[i]})
				successCount++
			}
		}

		err2 := tx.Model(users).Association("Favorites").Delete(videos)
		if err2 != nil {
			return err2
		}

		var userBatch []model.User
		err2 = tx.Model(users).FindInBatches(&userBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, user := range userBatch {
				err3 := bx.Model(user).Update("FavoritesCount", gorm.Expr("favorites_count-?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		var videoBatch []model.Video
		err2 = tx.Model(videos).FindInBatches(&videoBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, video := range videoBatch {
				var authorID uint
				err3 := bx.Model(video).Select("author_id").Find(&authorID).Error
				if err3 != nil {
					return err3
				}
				author := &model.User{ID: authorID}

				err3 = bx.Model(author).Update("FavoritedCount", gorm.Expr("favorited_count-?", 1)).Error
				if err3 != nil {
					return err3
				}

				err3 = bx.Model(video).Update("FavoritedCount", gorm.Expr("favorited_count-?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return 0 // 因事务回滚
	}

	return successCount
}

// 读取点赞(视频)列表 (select: Favorites.ID)
func ReadUserFavorites(ctx context.Context, id uint) (videoIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Favorites").Find(&videoIDs)
	if err != nil {
		return []uint{}, err
	}
	return videoIDs, nil
}

// 读取点赞(视频)数量
func CountUserFavorites(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.User{ID: id}).Select("FavoritesCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 读取获赞数量
func CountUserFavorited(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.User{ID: id}).Select("FavoritedCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 检查点赞关系
func CheckUserFavorites(ctx context.Context, id uint, videoID uint) (isFavorite bool) {
	DB := _db.WithContext(ctx)
	var results []model.Video
	err := DB.Model(&model.User{ID: id}).Select("id").Where("id=?", videoID).Limit(1).Association("Favorites").Find(&results)
	return err == nil && len(results) > 0
}

// 读取评论列表(num==-1时取消数量限制) (select: Comments.ID)
func ReadUserComments(ctx context.Context, id uint, createdAt int64, forward bool, num int) (commentIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.User{ID: id}).Select("id").Where("created_at>?", stop).Order("created_at").Limit(num).Association("Comments").Find(&commentIDs)
	} else {
		err = DB.Model(&model.User{ID: id}).Select("id").Where("created_at<?", stop).Order("created_at desc").Limit(num).Association("Comments").Find(&commentIDs)
	}
	if err != nil {
		return []uint{}, err
	}
	return commentIDs, nil
}

// 读取评论数量
func CountUserComments(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.User{ID: id}).Select("CommentsCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 检查评论所属
func CheckUserComments(ctx context.Context, id uint, commentID uint) (isIts bool) {
	DB := _db.WithContext(ctx)
	var results []model.Comment
	err := DB.Model(&model.User{ID: id}).Select("id").Where("id=?", commentID).Limit(1).Association("Comments").Find(&results)
	return err == nil && len(results) > 0
}

// 创建关注关系
func CreateUserFollows(ctx context.Context, id uint, followID uint) (err error) {
	if id == followID {
		return ErrorSelfFollow // 默认禁止自己关注自己
	}

	DB := _db.WithContext(ctx)
	return DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		user := &model.User{ID: id}
		follow := &model.User{ID: followID}

		var results []model.User
		err2 := tx.Model(user).Select("id").Where("id=?", followID).Limit(1).Association("Follows").Find(&results)
		if err2 != nil {
			return err2
		}
		if len(results) > 0 { // 不允许重复创建
			return ErrorRecordExists
		}

		err2 = tx.Model(user).Association("Follows").Append(follow)
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(user).Update("FollowsCount", gorm.Expr("follows_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(follow).Update("FollowersCount", gorm.Expr("followers_count+?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
}

// 创建关注关系(批量)
func CreateUserFollowsBatch(ctx context.Context, ids []uint, followIDs []uint) (successCount int64) {
	successCount = 0
	if len(ids) != len(followIDs) {
		return 0
	}

	DB := _db.WithContext(ctx)
	err := DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		users := make([]model.User, 0, len(ids))
		follows := make([]model.User, 0, len(followIDs))
		for i, id := range ids {
			var results []model.User
			if tx.Model(&model.User{ID: id}).Select("id").Where("id=?", followIDs[i]).Limit(1).Association("Follows").Find(&results) == nil && len(results) == 0 { // 不允许重复创建
				users = append(users, model.User{ID: id})
				follows = append(follows, model.User{ID: followIDs[i]})
				successCount++
			}
		}

		err2 := tx.Model(users).Association("Follows").Append(follows)
		if err2 != nil {
			return err2
		}

		var userBatch []model.User
		err2 = tx.Model(users).FindInBatches(&userBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, user := range userBatch {
				err3 := bx.Model(user).Update("FollowsCount", gorm.Expr("follows_count+?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		var followBatch []model.User
		err2 = tx.Model(follows).FindInBatches(&followBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, follow := range followBatch {
				err3 := bx.Model(follow).Update("FollowersCount", gorm.Expr("followers_count+?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return 0 // 因事务回滚
	}

	return successCount
}

// 删除关注关系
func DeleteUserFollows(ctx context.Context, id uint, followID uint) (err error) {
	DB := _db.WithContext(ctx)
	return DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		user := &model.User{ID: id}
		follow := &model.User{ID: followID}

		var results []model.Video
		err2 := tx.Model(user).Select("id").Where("id=?", followID).Limit(1).Association("Follows").Find(&results)
		if err2 != nil {
			return err2
		}
		if len(results) == 0 { // 不允许凭空删除
			return ErrorRecordNotExists
		}

		err2 = tx.Model(user).Association("Follows").Delete(follow)
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(user).Update("FollowsCount", gorm.Expr("follows_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		err2 = tx.Model(follow).Update("FollowersCount", gorm.Expr("followers_count-?", 1)).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
}

// 删除关注关系(批量)
func DeleteUserFollowsBatch(ctx context.Context, ids []uint, followIDs []uint) (successCount int64) {
	successCount = 0
	if len(ids) != len(followIDs) {
		return 0
	}

	DB := _db.WithContext(ctx)
	err := DB.Transaction(func(tx *gorm.DB) error { // 使用事务
		users := make([]model.User, 0, len(ids))
		follows := make([]model.User, 0, len(followIDs))
		for i, id := range ids {
			var results []model.User
			if tx.Model(&model.User{ID: id}).Select("id").Where("id=?", followIDs[i]).Limit(1).Association("Follows").Find(&results) == nil && len(results) > 0 { // 不允许凭空删除
				users = append(users, model.User{ID: id})
				follows = append(follows, model.User{ID: followIDs[i]})
				successCount++
			}
		}

		err2 := tx.Model(users).Association("Follows").Delete(follows)
		if err2 != nil {
			return err2
		}

		var userBatch []model.User
		err2 = tx.Model(users).FindInBatches(&userBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, user := range userBatch {
				err3 := bx.Model(user).Update("FollowsCount", gorm.Expr("follows_count-?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		var followBatch []model.User
		err2 = tx.Model(follows).FindInBatches(&followBatch, batchNum, func(bx *gorm.DB, batch int) error {
			for _, follow := range followBatch {
				err3 := bx.Model(follow).Update("FollowersCount", gorm.Expr("followers_count-?", 1)).Error
				if err3 != nil {
					return err3
				}
			}
			return nil
		}).Error
		if err2 != nil {
			return err2
		}

		return nil
	})
	if err != nil {
		return 0 // 因事务回滚
	}

	return successCount
}

// 读取关注(用户)列表 (select: Follows.ID)
func ReadUserFollows(ctx context.Context, id uint) (userIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Follows").Find(&userIDs)
	if err != nil {
		return []uint{}, err
	}
	return userIDs, nil
}

// 读取关注(用户)数量
func CountUserFollows(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.User{ID: id}).Select("FollowsCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 读取粉丝(用户)列表 (select: Followers.ID)
func ReadUserFollowers(ctx context.Context, id uint) (userIDs []uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Followers").Find(&userIDs)
	if err != nil {
		return []uint{}, err
	}
	return userIDs, nil
}

// 读取粉丝(用户)数量
func CountUserFollowers(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	err := DB.Model(&model.User{ID: id}).Select("FollowersCount").Scan(&count).Error
	if err != nil {
		return -1 // 出错
	}
	return count
}

// 检查关注关系
func CheckUserFollows(ctx context.Context, id uint, followID uint) (isFollowing bool) {
	if id == followID {
		return false // 默认自己不关注自己
	}

	DB := _db.WithContext(ctx)
	var results []model.User
	err := DB.Model(&model.User{ID: id}).Select("id").Where("id=?", followID).Limit(1).Association("Follows").Find(&results)
	return err == nil && len(results) > 0
}
