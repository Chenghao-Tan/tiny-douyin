package db

import (
	"douyin/repo/internal/db/model"

	"context"
	"errors"

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
func ReadUserWorks(ctx context.Context, id uint) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Works").Find(&videos)
	if err != nil {
		return videos, err
	}
	return videos, nil
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

// 读取点赞(视频)列表 (select: Favorites.ID)
func ReadUserFavorites(ctx context.Context, id uint) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Favorites").Find(&videos)
	if err != nil {
		return videos, err
	}
	return videos, nil
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

// 读取评论列表 (select: Comments.ID)
func ReadUserComments(ctx context.Context, id uint) (comments []model.Comment, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Comments").Find(&comments)
	if err != nil {
		return comments, err
	}
	return comments, nil
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

// 读取关注(用户)列表 (select: Follows.ID)
func ReadUserFollows(ctx context.Context, id uint) (users []model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Follows").Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
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
func ReadUserFollowers(ctx context.Context, id uint) (users []model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Followers").Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
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

// 读取消息列表 (select: Messages.ID)
func ReadUserMessages(ctx context.Context, id uint) (messages []model.Message, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{ID: id}).Select("id").Association("Messages").Find(&messages)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

// 计算消息数量
func CountUserMessages(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{ID: id}).Association("Messages").Count()
}
