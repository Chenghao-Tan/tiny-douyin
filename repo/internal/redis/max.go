package redis

import (
	"context"
)

const prefixUserMaxID = "user:max"       // 直接读取
const prefixVideoMaxID = "video:max"     // 直接读取
const prefixCommentMaxID = "comment:max" // 直接读取
const prefixMessageMaxID = "message:max" // 直接读取

// 设置用户主键最大值
func SetUserMaxID(ctx context.Context, maxID uint) (err error) {
	key := prefixUserMaxID
	return _redis.Set(ctx, key, maxID, 0).Err()
}

// 读取用户主键最大值
func GetUserMaxID(ctx context.Context) (maxID uint, err error) {
	key := prefixUserMaxID
	value, err := _redis.Get(ctx, key).Uint64()
	return uint(value), err
}

// 自增用户主键最大值
func IncrUserMaxID(ctx context.Context) (err error) {
	key := prefixUserMaxID
	return _redis.Incr(ctx, key).Err()
}

// 设置视频主键最大值
func SetVideoMaxID(ctx context.Context, maxID uint) (err error) {
	key := prefixVideoMaxID
	return _redis.Set(ctx, key, maxID, 0).Err()
}

// 读取视频主键最大值
func GetVideoMaxID(ctx context.Context) (maxID uint, err error) {
	key := prefixVideoMaxID
	value, err := _redis.Get(ctx, key).Uint64()
	return uint(value), err
}

// 自增视频主键最大值
func IncrVideoMaxID(ctx context.Context) (err error) {
	key := prefixVideoMaxID
	return _redis.Incr(ctx, key).Err()
}

// 设置评论主键最大值
func SetCommentMaxID(ctx context.Context, maxID uint) (err error) {
	key := prefixCommentMaxID
	return _redis.Set(ctx, key, maxID, 0).Err()
}

// 读取评论主键最大值
func GetCommentMaxID(ctx context.Context) (maxID uint, err error) {
	key := prefixCommentMaxID
	value, err := _redis.Get(ctx, key).Uint64()
	return uint(value), err
}

// 自增评论主键最大值
func IncrCommentMaxID(ctx context.Context) (err error) {
	key := prefixCommentMaxID
	return _redis.Incr(ctx, key).Err()
}

// 设置消息主键最大值
func SetMessageMaxID(ctx context.Context, maxID uint) (err error) {
	key := prefixMessageMaxID
	return _redis.Set(ctx, key, maxID, 0).Err()
}

// 读取消息主键最大值
func GetMessageMaxID(ctx context.Context) (maxID uint, err error) {
	key := prefixMessageMaxID
	value, err := _redis.Get(ctx, key).Uint64()
	return uint(value), err
}

// 自增消息主键最大值
func IncrMessageMaxID(ctx context.Context) (err error) {
	key := prefixMessageMaxID
	return _redis.Incr(ctx, key).Err()
}
