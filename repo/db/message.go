package db

import (
	"douyin/repo/db/model"

	"context"
	"time"

	"gorm.io/gorm/clause"
)

// 发送消息
func CreateMessage(ctx context.Context, from_user_id uint, to_user_id uint, content string) (message *model.Message, err error) {
	DB := _db.WithContext(ctx)
	message = &model.Message{Content: content, FromUserID: from_user_id, ToUserID: to_user_id}
	err = DB.Model(&model.Message{}).Create(message).Error
	return message, err
}

// 根据发送者ID和接收者ID获取消息列表
func FindMessagesBy_From_To_ID(ctx context.Context, from_user_id uint, to_user_id uint, createdAt int64, forward bool, num int) (messages []model.Message, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.Message{}).Where("created_at>?", stop).Where("from_user_id=? AND to_user_id=?", from_user_id, to_user_id).Order("created_at").Limit(num).Preload(clause.Associations).Find(&messages).Error
	} else {
		err = DB.Model(&model.Message{}).Where("created_at<?", stop).Where("from_user_id=? AND to_user_id=?", from_user_id, to_user_id).Order("created_at desc").Limit(num).Preload(clause.Associations).Find(&messages).Error
	}
	return messages, err
}
