package db

import (
	"douyin/repo/db/model"

	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrorSelfFollow = errors.New("禁止自己关注自己")

// 关注
func CreateFollow(ctx context.Context, user_id uint, follow_id uint) (err error) {
	if user_id == follow_id {
		return ErrorSelfFollow
	}

	DB := _db.WithContext(ctx)
	user := &model.User{Model: gorm.Model{ID: user_id}}
	follow := &model.User{Model: gorm.Model{ID: follow_id}}
	return DB.Model(user).Association("Follows").Append(follow)
}

// 取消关注
func DeleteFollow(ctx context.Context, user_id uint, follow_id uint) (err error) {
	DB := _db.WithContext(ctx)
	user := &model.User{Model: gorm.Model{ID: user_id}}
	follow := &model.User{Model: gorm.Model{ID: follow_id}}
	return DB.Model(user).Association("Follows").Delete(follow)
}

// 检查是否关注
func CheckFollow(ctx context.Context, user_id uint, follow_id uint) (is_follower bool) {
	if user_id == follow_id {
		return false // 默认自己不关注自己
	}

	DB := _db.WithContext(ctx)
	user := &model.User{Model: gorm.Model{ID: user_id}}
	return DB.Model(user).Where("id=?", follow_id).Association("Follows").Count() > 0
}
