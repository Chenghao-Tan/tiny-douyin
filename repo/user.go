package repo

import (
	"douyin/repo/db"
	"douyin/repo/redis"

	"context"
	"strconv"
)

// 创建点赞关系
func CreateUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("fav:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(videoID), 10) + ":1")

	video, err := db.ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID
	if err != nil {
		return err
	}
	return redis.SetUserFavorites(ctx, id, videoID, video.AuthorID, true, maxSyncDelay)
}

// 删除点赞关系
func DeleteUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("fav:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(videoID), 10) + ":0")

	video, err := db.ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID
	if err != nil {
		return err
	}
	return redis.SetUserFavorites(ctx, id, videoID, video.AuthorID, false, maxSyncDelay)
}

// 读取点赞(视频)数量
func CountUserFavorites(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserFavoritesCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserFavoritesCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserFavorites(ctx, id)
		if record >= 0 {
			_ = redis.SetUserFavoritesCount(ctx, id, record, cacheExpiration)
			return record
		}
	}
	return 0 // 当出现错误
}

// 读取获赞数量
func CountUserFavorited(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserFavoritedCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserFavoritedCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserFavorited(ctx, id)
		if record >= 0 {
			_ = redis.SetUserFavoritedCount(ctx, id, record, cacheExpiration)
			return record
		}
	}
	return 0 // 当出现错误
}

// 检查点赞关系
func CheckUserFavorites(ctx context.Context, id uint, videoID uint) (isFavorite bool) {
	isFavorite, err := redis.GetUserFavorites(ctx, id, videoID, distrustProbability)
	if err == nil { // 命中缓存
		return isFavorite
	}
	if err == redis.ErrorRedisNil { // 启动同步
		record := db.CheckUserFavorites(ctx, id, videoID)
		_ = redis.SetUserFavoritesBit(ctx, id, videoID, record) // 立即修正缓存主记录
		return record
	}
	return false // 当出现错误
}
