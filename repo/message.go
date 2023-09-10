package repo

import (
	"douyin/repo/internal/db"
	"douyin/repo/internal/db/model"
	"douyin/repo/internal/redis"

	"context"
)

// 获取消息主键最大值
func MaxMessageID(ctx context.Context) (id uint, err error) {
	return redis.GetMessageMaxID(ctx)
}

// 创建消息 //TODO
func CreateMessage(ctx context.Context, fromUserID uint, toUserID uint, content string) (message *model.Message, err error) {
	message, err = db.CreateMessage(ctx, fromUserID, toUserID, content)
	if err != nil {
		return nil, err
	}
	_ = redis.IncrMessageMaxID(ctx)
	return message, nil
}

// 根据聊天双方ID和创建时间查找消息列表(num==-1时取消数量限制) (select: ID, CreatedAt, Content) //TODO
func FindMessagesByCreatedAt(ctx context.Context, User1ID uint, User2ID uint, createdAt int64, forward bool, num int) (messages []model.Message, err error) {
	return db.FindMessagesByCreatedAt(ctx, User1ID, User2ID, createdAt, forward, num)
}
