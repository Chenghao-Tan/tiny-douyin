package repo

import (
	"douyin/repo/db"
	"douyin/repo/db/model"
	"douyin/repo/redis"

	"context"
)

// 读取视频基本信息 (select: *)
func ReadVideoBasics(ctx context.Context, id uint) (video *model.Video, err error) {
	video, err = redis.GetVideoBasics(ctx, id)
	if err == nil { // 命中缓存
		return video, nil
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetVideoBasics(ctx, id, &model.Video{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadVideoBasics(ctx, id)
		if err != nil {
			_ = redis.SetVideoBasics(ctx, id, record, cacheExpiration)
			return record, nil
		}
	}
	return nil, err // 当出现错误
}

// 读取点赞(用户)数量
func CountVideoFavorited(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetVideoFavoritedCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetVideoFavoritedCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountVideoFavorited(ctx, id)
		if record >= 0 {
			_ = redis.SetVideoFavoritedCount(ctx, id, record, cacheExpiration)
			return record
		}
	}
	return 0 // 当出现错误
}
