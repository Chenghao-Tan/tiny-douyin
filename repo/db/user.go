package db

import (
	"douyin/repo/db/model"

	"context"
	"errors"

	"gorm.io/gorm"
)

// 自定义错误类型
var ErrorSelfFollow = errors.New("禁止自己关注自己")

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

// 读取用户基本信息 (select: *)
func ReadUserBasics(ctx context.Context, id uint) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("id=?", id).First(user).Error
	return user, err
}

// 根据用户名读取用户基本信息 (select: *)
func ReadUserBasicsByUsername(ctx context.Context, username string) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("username=?", username).First(user).Error
	return user, err
}

// 读取作品(视频)列表 (select: Works.ID)
func ReadUserWorks(ctx context.Context, id uint) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Association("Works").Find(&videos)
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// 读取作品(视频)数量
func CountUserWorks(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Works").Count()
}

// 读取点赞(视频)列表 (select: Favorites.ID)
func ReadUserFavorites(ctx context.Context, id uint) (videos []model.Video, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Association("Favorites").Find(&videos)
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// 读取点赞(视频)数量
func CountUserFavorites(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Favorites").Count()
}

// 创建点赞关系
func CreateUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	DB := _db.WithContext(ctx)
	video := &model.Video{Model: gorm.Model{ID: videoID}}
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Favorites").Append(video)
}

// 删除点赞关系
func DeleteUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	DB := _db.WithContext(ctx)
	video := &model.Video{Model: gorm.Model{ID: videoID}}
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Favorites").Delete(video)
}

// 检查点赞关系
func CheckUserFavorites(ctx context.Context, id uint, videoID uint) (isFavorite bool) {
	DB := _db.WithContext(ctx)
	var results []model.Video
	err := DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Where("id=?", videoID).Limit(1).Association("Favorites").Find(&results)
	return err == nil && len(results) > 0
}

// 读取获赞数量
func CountUserFavorited(ctx context.Context, id uint) (count int64) {
	works, err := ReadUserWorks(ctx, id)
	if err != nil {
		return 0 //TODO (可为-1)
	}
	count = 0
	for _, video := range works {
		count += CountVideoFavorited(ctx, video.ID)
	}
	return count
}

// 读取评论列表 (select: Comments.ID)
func ReadUserComments(ctx context.Context, id uint) (comments []model.Comment, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Association("Comments").Find(&comments)
	if err != nil {
		return comments, err
	}
	return comments, nil
}

// 读取评论数量
func CountUserComments(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Comments").Count()
}

// 读取关注(用户)列表 (select: Follows.ID)
func ReadUserFollows(ctx context.Context, id uint) (users []model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Association("Follows").Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

// 读取关注(用户)数量
func CountUserFollows(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Follows").Count()
}

// 创建关注关系
func CreateUserFollows(ctx context.Context, id uint, followID uint) (err error) {
	if id == followID {
		return ErrorSelfFollow // 默认禁止自己关注自己
	}

	DB := _db.WithContext(ctx)
	follow := &model.User{Model: gorm.Model{ID: followID}}
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Follows").Append(follow)
}

// 删除关注关系
func DeleteUserFollows(ctx context.Context, id uint, followID uint) (err error) {
	DB := _db.WithContext(ctx)
	follow := &model.User{Model: gorm.Model{ID: followID}}
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Follows").Delete(follow)
}

// 检查关注关系
func CheckUserFollows(ctx context.Context, id uint, followID uint) (isFollowing bool) {
	if id == followID {
		return false // 默认自己不关注自己
	}

	DB := _db.WithContext(ctx)
	var results []model.User
	err := DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Where("id=?", followID).Limit(1).Association("Follows").Find(&results)
	return err == nil && len(results) > 0
}

// 读取粉丝(用户)列表 (select: Followers.ID)
func ReadUserFollowers(ctx context.Context, id uint) (users []model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Association("Followers").Find(&users)
	if err != nil {
		return users, err
	}
	return users, nil
}

// 读取粉丝(用户)数量
func CountUserFollowers(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Followers").Count()
}

// 读取消息列表 (select: Messages.ID)
func ReadUserMessages(ctx context.Context, id uint) (messages []model.Message, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{Model: gorm.Model{ID: id}}).Select("id").Association("Messages").Find(&messages)
	if err != nil {
		return messages, err
	}
	return messages, nil
}

// 读取消息数量
func CountUserMessages(ctx context.Context, id uint) (count int64) {
	DB := _db.WithContext(ctx)
	return DB.Model(&model.User{Model: gorm.Model{ID: id}}).Association("Messages").Count()
}
