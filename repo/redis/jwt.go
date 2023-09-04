package redis

import (
	"context"
	"strconv"
	"time"
)

const prefixUserJWT = "user:jwt:" // 后接三十六进制userID (节约key长度)

// 设置用户token
func SetUserJWT(ctx context.Context, userID uint, token string, autoLogout time.Duration) (err error) {
	key := prefixUserJWT + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, token, autoLogout).Err()
}

// 读取用户token
func GetUserJWT(ctx context.Context, userID uint) (token string, err error) {
	key := prefixUserJWT + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Result()
}

// 设置用户token过期时间
func ExpireUserJWT(ctx context.Context, userID uint, autoLogout time.Duration) (err error) {
	key := prefixUserJWT + strconv.FormatUint(uint64(userID), 36)
	return _redis.Expire(ctx, key, autoLogout).Err()
}
