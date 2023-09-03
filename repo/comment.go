package repo

import (
	"douyin/repo/db"
	"douyin/repo/db/model"
	"douyin/repo/redis"

	"context"
)

// 读取评论基本信息 (select: *)
func ReadCommentBasics(ctx context.Context, id uint) (comment *model.Comment, err error) {
	comment, err = redis.GetCommentBasics(ctx, id)
	if err == nil { // 命中缓存
		return comment, nil
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetCommentBasics(ctx, id, &model.Comment{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadCommentBasics(ctx, id)
		if err != nil {
			_ = redis.SetCommentBasics(ctx, id, record, cacheExpiration)
			return record, nil
		}
	}
	return nil, err // 当出现错误
}
