package redis

import (
	"context"
	"strconv"
	"time"
)

const prefixUserWorks = "user:wrk:"                     // 暂只用于构建其他前缀
const prefixUserWorksCount = prefixUserWorks + "count:" // 后接三十六进制userID (节约key长度)

// 设置用户作品数
func SetUserWorksCount(ctx context.Context, userID uint, count int64, expiration time.Duration) (err error) {
	key := prefixUserWorksCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取用户作品数
func GetUserWorksCount(ctx context.Context, userID uint) (count int64, err error) {
	key := prefixUserWorksCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Int64()
}

// 删除用户作品数
func DelUserWorksCount(ctx context.Context, userID uint, maxWriteTime time.Duration) (err error) {
	key := prefixUserWorksCount + strconv.FormatUint(uint64(userID), 36)
	err = _redis.Del(ctx, key).Err()

	// 缓存双删
	go func() {
		time.Sleep(maxWriteTime)

		_ = _redis.Del(ctx, key).Err()
	}()

	return err
}
