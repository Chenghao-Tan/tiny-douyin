package redis

import (
	"douyin/repo/db/model"

	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const prefixUserBasics = "user:bsc:"       // 后接三十六进制userID (节约key长度)
const prefixVideoBasics = "video:bsc:"     // 后接三十六进制videoID (节约key长度)
const prefixCommentBasics = "comment:bsc:" // 后接三十六进制commentID (节约key长度)

// 设置用户基本信息
func SetUserBasics(ctx context.Context, userID uint, user *model.User, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixUserBasics + strconv.FormatUint(uint64(userID), 36)

		pipe.HSet(ctx, key, user)
		pipe.Expire(ctx, key, expiration)

		return nil
	})
	return err
}

// 读取用户基本信息
func GetUserBasics(ctx context.Context, userID uint) (user *model.User, err error) {
	key := prefixUserBasics + strconv.FormatUint(uint64(userID), 36)
	user = &model.User{}
	err = _redis.HGetAll(ctx, key).Scan(user)
	return user, err
}

// 设置视频基本信息
func SetVideoBasics(ctx context.Context, videoID uint, video *model.Video, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixVideoBasics + strconv.FormatUint(uint64(videoID), 36)

		pipe.HSet(ctx, key, video)
		pipe.Expire(ctx, key, expiration)

		return nil
	})
	return err
}

// 读取视频基本信息
func GetVideoBasics(ctx context.Context, videoID uint) (video *model.Video, err error) {
	key := prefixVideoBasics + strconv.FormatUint(uint64(videoID), 36)
	video = &model.Video{}
	err = _redis.HGetAll(ctx, key).Scan(video)
	return video, err
}

// 设置评论基本信息
func SetCommentBasics(ctx context.Context, commentID uint, comment *model.Comment, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixCommentBasics + strconv.FormatUint(uint64(commentID), 36)

		pipe.HSet(ctx, key, comment)
		pipe.Expire(ctx, key, expiration)

		return nil
	})
	return err
}

// 读取评论基本信息
func GetCommentBasics(ctx context.Context, commentID uint) (comment *model.Comment, err error) {
	key := prefixCommentBasics + strconv.FormatUint(uint64(commentID), 36)
	comment = &model.Comment{}
	err = _redis.HGetAll(ctx, key).Scan(comment)
	return comment, err
}
