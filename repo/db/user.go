package db

import (
	"douyin/repo/db/model"

	"context"

	"gorm.io/gorm/clause"
)

// 创建User
func CreateUser(ctx context.Context, userInfo *model.User) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.User{}).Create(userInfo).Error
	return userInfo, err
}

// 根据用户ID找到用户
func FindUserByID(ctx context.Context, id uint) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("id=?", id).Preload(clause.Associations).Preload("Works.Favorited").Preload("Works.Comments").Preload("Favorites.Favorited").Preload("Favorites.Comments").First(user).Error
	return user, err
}

// 根据用户名找到用户
func FindUserByUsername(ctx context.Context, username string) (user *model.User, err error) {
	DB := _db.WithContext(ctx)
	user = &model.User{}
	err = DB.Model(&model.User{}).Where("username=?", username).Preload(clause.Associations).Preload("Works.Favorited").Preload("Works.Comments").Preload("Favorites.Favorited").Preload("Favorites.Comments").First(user).Error
	return user, err
}
