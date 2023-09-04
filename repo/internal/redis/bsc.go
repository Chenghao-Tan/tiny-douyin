package redis

import (
	"douyin/repo/internal/db/model"

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
		pipe.Expire(ctx, key, randomExpiration(expiration))

		return nil
	})
	return err
}

// 读取用户基本信息
func GetUserBasics(ctx context.Context, userID uint) (user *model.User, err error) {
	key := prefixUserBasics + strconv.FormatUint(uint64(userID), 36)

	var exists int64
	var cmd *redis.MapStringStringCmd
	for i := 0; i < maxRetries; i++ {
		err = _redis.Watch(ctx, func(tx *redis.Tx) error { // 使用乐观锁
			var err2 error
			exists, err2 = tx.Exists(ctx, key).Result()
			if err2 != nil && err2 != ErrorRedisNil {
				return err2
			}

			_, err2 = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
				cmd = pipe.HGetAll(ctx, key) // 返回的错误只会表示异常, 不会为ErrorRedisNil
				return nil
			})
			return err2
		}, key)
		if err == nil { // 乐观锁成功
			break
		} else if err == ErrorRedisTxFailed { // 乐观锁失败, 但无其他异常
			continue
		} else { // 出现其他异常
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	// 乐观锁成功, 结果可信
	if exists == 0 {
		return nil, ErrorRedisNil
	}
	user = &model.User{}
	err = cmd.Scan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 设置视频基本信息
func SetVideoBasics(ctx context.Context, videoID uint, video *model.Video, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixVideoBasics + strconv.FormatUint(uint64(videoID), 36)

		pipe.HSet(ctx, key, video)
		pipe.Expire(ctx, key, randomExpiration(expiration))

		return nil
	})
	return err
}

// 读取视频基本信息
func GetVideoBasics(ctx context.Context, videoID uint) (video *model.Video, err error) {
	key := prefixVideoBasics + strconv.FormatUint(uint64(videoID), 36)

	var exists int64
	var cmd *redis.MapStringStringCmd
	for i := 0; i < maxRetries; i++ {
		err = _redis.Watch(ctx, func(tx *redis.Tx) error { // 使用乐观锁
			var err2 error
			exists, err2 = tx.Exists(ctx, key).Result()
			if err2 != nil && err2 != ErrorRedisNil {
				return err2
			}

			_, err2 = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
				cmd = pipe.HGetAll(ctx, key) // 返回的错误只会表示异常, 不会为ErrorRedisNil
				return nil
			})
			return err2
		}, key)
		if err == nil { // 乐观锁成功
			break
		} else if err == ErrorRedisTxFailed { // 乐观锁失败, 但无其他异常
			continue
		} else { // 出现其他异常
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	// 乐观锁成功, 结果可信
	if exists == 0 {
		return nil, ErrorRedisNil
	}
	video = &model.Video{}
	err = cmd.Scan(video)
	if err != nil {
		return nil, err
	}
	return video, nil
}

// 设置评论基本信息
func SetCommentBasics(ctx context.Context, commentID uint, comment *model.Comment, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixCommentBasics + strconv.FormatUint(uint64(commentID), 36)

		pipe.HSet(ctx, key, comment)
		pipe.Expire(ctx, key, randomExpiration(expiration))

		return nil
	})
	return err
}

// 读取评论基本信息
func GetCommentBasics(ctx context.Context, commentID uint) (comment *model.Comment, err error) {
	key := prefixCommentBasics + strconv.FormatUint(uint64(commentID), 36)

	var exists int64
	var cmd *redis.MapStringStringCmd
	for i := 0; i < maxRetries; i++ {
		err = _redis.Watch(ctx, func(tx *redis.Tx) error { // 使用乐观锁
			var err2 error
			exists, err2 = tx.Exists(ctx, key).Result()
			if err2 != nil && err2 != ErrorRedisNil {
				return err2
			}

			_, err2 = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
				cmd = pipe.HGetAll(ctx, key) // 返回的错误只会表示异常, 不会为ErrorRedisNil
				return nil
			})
			return err2
		}, key)
		if err == nil { // 乐观锁成功
			break
		} else if err == ErrorRedisTxFailed { // 乐观锁失败, 但无其他异常
			continue
		} else { // 出现其他异常
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	// 乐观锁成功, 结果可信
	if exists == 0 {
		return nil, ErrorRedisNil
	}
	comment = &model.Comment{}
	err = cmd.Scan(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
