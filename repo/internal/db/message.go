package db

import (
	"douyin/repo/internal/db/model"

	"context"
	"time"
)

// 获取消息主键最大值
func MaxMessageID(ctx context.Context) (id uint, err error) {
	DB := _db.WithContext(ctx)
	err = DB.Model(&model.Message{}).Select("MAX(id)").Scan(&id).Error
	return id, err
}

// 创建消息
func CreateMessage(ctx context.Context, fromUserID uint, toUserID uint, content string) (message *model.Message, err error) {
	DB := _db.WithContext(ctx)
	message = &model.Message{Content: content, FromUserID: fromUserID, ToUserID: toUserID}
	err = DB.Model(&model.Message{}).Create(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

// 根据聊天双方ID和创建时间查找消息列表(num==-1时取消数量限制) (select: *)
func FindMessagesByCreatedAt(ctx context.Context, User1ID uint, User2ID uint, createdAt int64, forward bool, num int) (messages []model.Message, err error) {
	DB := _db.WithContext(ctx)
	stop := time.Unix(createdAt, 0)
	if forward {
		err = DB.Model(&model.Message{}).Where(DB.Where("from_user_id=? AND to_user_id=?", User1ID, User2ID).Or("from_user_id=? AND to_user_id=?", User2ID, User1ID)).Where("created_at>?", stop).Order("created_at").Limit(num).Find(&messages).Error
	} else {
		err = DB.Model(&model.Message{}).Where(DB.Where("from_user_id=? AND to_user_id=?", User1ID, User2ID).Or("from_user_id=? AND to_user_id=?", User2ID, User1ID)).Where("created_at<?", stop).Order("created_at desc").Limit(num).Find(&messages).Error
	}
	if err != nil {
		return messages, err
	}
	return messages, err
}
