package repo

import (
	"douyin/repo/internal/redis"

	"context"
	"time"
)

// 设置用户token
func SetUserJWT(ctx context.Context, userID uint, token string, autoLogout time.Duration) (err error) {
	return redis.SetUserJWT(ctx, userID, token, autoLogout)
}

// 读取用户token
func GetUserJWT(ctx context.Context, userID uint) (token string, err error) {
	return redis.GetUserJWT(ctx, userID)
}

// 设置用户token过期时间
func ExpireUserJWT(ctx context.Context, userID uint, autoLogout time.Duration) (err error) {
	return redis.ExpireUserJWT(ctx, userID, autoLogout)
}
