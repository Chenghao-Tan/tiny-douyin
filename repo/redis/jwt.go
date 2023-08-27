package redis

import (
	"context"
	"strconv"
	"time"
)

const prefixJWT = "user:jwt:" // 后接三十六进制userID (节约key长度)

// 设置用户token
func SetJWT(ctx context.Context, userID uint, token string, expiration time.Duration) (err error) {
	key := prefixJWT + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, token, expiration).Err()
}

// 读取用户token
func GetJWT(ctx context.Context, userID uint) (token string, err error) {
	key := prefixJWT + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Result()
}

// 设置用户token过期时间
func ExpireJWT(ctx context.Context, userID uint, expiration time.Duration) (err error) {
	key := prefixJWT + strconv.FormatUint(uint64(userID), 36)
	return _redis.Expire(ctx, key, expiration).Err()
}

// 检查用户token
func CheckJWT(ctx context.Context, userID uint, token string) (isValid bool) {
	record, err := GetJWT(ctx, userID)
	if err != nil || token != record {
		return false
	}
	return true
}
