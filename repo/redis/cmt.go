package redis

import (
	"context"
	"strconv"
	"time"
)

const prefixUserComments = "user:cmt:"                          // 暂只用于构建其他前缀
const prefixUserCommentsCount = prefixUserComments + "count:"   // 后接三十六进制userID (节约key长度)
const prefixVideoComments = "video:cmt:"                        // 暂只用于构建其他前缀
const prefixVideoCommentsCount = prefixVideoComments + "count:" // 后接三十六进制videoID (节约key长度)

// 设置用户评论数
func SetUserCommentsCount(ctx context.Context, userID uint, count int64, expiration time.Duration) (err error) {
	key := prefixUserCommentsCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取用户评论数
func GetUserCommentsCount(ctx context.Context, userID uint) (count int64, err error) {
	key := prefixUserCommentsCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Int64()
}

// 删除用户评论数
func DelUserCommentsCount(ctx context.Context, userID uint, maxWriteTime time.Duration) (err error) {
	key := prefixUserCommentsCount + strconv.FormatUint(uint64(userID), 36)
	err = _redis.Del(ctx, key).Err()

	// 缓存双删
	go func() {
		time.Sleep(maxWriteTime)

		_ = _redis.Del(ctx, key).Err()
	}()

	return err
}

// 设置视频评论数
func SetVideoCommentsCount(ctx context.Context, videoID uint, count int64, expiration time.Duration) (err error) {
	key := prefixVideoCommentsCount + strconv.FormatUint(uint64(videoID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取视频评论数
func GetVideoCommentsCount(ctx context.Context, videoID uint) (count int64, err error) {
	key := prefixVideoCommentsCount + strconv.FormatUint(uint64(videoID), 36)
	return _redis.Get(ctx, key).Int64()
}

// 删除视频评论数
func DelVideoCommentsCount(ctx context.Context, videoID uint, maxWriteTime time.Duration) (err error) {
	key := prefixVideoCommentsCount + strconv.FormatUint(uint64(videoID), 36)
	err = _redis.Del(ctx, key).Err()

	// 缓存双删
	go func() {
		time.Sleep(maxWriteTime)

		_ = _redis.Del(ctx, key).Err()
	}()

	return err
}
