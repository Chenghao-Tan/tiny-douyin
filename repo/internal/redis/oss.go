package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const prefixUserOSS = "user:oss:"                               // 暂只用于构建其他前缀
const prefixUserAvatarURL = prefixUserOSS + "avatar:"           // 后接objectID
const prefixUserBackgroundImageURL = prefixUserOSS + "bgimage:" // 后接objectID
const prefixVideoOSS = "video:oss:"                             // 暂只用于构建其他前缀
const prefixVideoURL = prefixVideoOSS                           // 后接objectID

// 设置头像对象外链
func SetUserAvatarURL(ctx context.Context, objectID string, avatarURL string, urlExpiration time.Duration) (err error) {
	key := prefixUserAvatarURL + objectID
	return _redis.SetEx(ctx, key, avatarURL, urlExpiration).Err()
}

// 读取头像对象外链
func GetUserAvatarURL(ctx context.Context, objectID string) (avatarURL string, err error) {
	key := prefixUserAvatarURL + objectID
	return _redis.Get(ctx, key).Result()
}

// 设置个人页背景图对象外链
func SetUserBackgroundImageURL(ctx context.Context, objectID string, backgroundImageURL string, urlExpiration time.Duration) (err error) {
	key := prefixUserBackgroundImageURL + objectID
	return _redis.SetEx(ctx, key, backgroundImageURL, urlExpiration).Err()
}

// 读取个人页背景图对象外链
func GetUserBackgroundImageURL(ctx context.Context, objectID string) (backgroundImageURL string, err error) {
	key := prefixUserBackgroundImageURL + objectID
	return _redis.Get(ctx, key).Result()
}

// 视频及封面对象结构体
type videoOSS struct {
	VideoURL string `redis:"video"`
	CoverURL string `redis:"cover"`
}

// 设置视频对象及封面对象外链
func SetVideoURL(ctx context.Context, objectID string, videoURL string, coverURL string, urlExpiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixVideoURL + objectID
		value := &videoOSS{VideoURL: videoURL, CoverURL: coverURL}

		pipe.HSet(ctx, key, value)
		pipe.Expire(ctx, key, urlExpiration)

		return nil
	})
	return err
}

// 读取视频对象及封面对象外链
func GetVideoURL(ctx context.Context, objectID string) (videoURL string, coverURL string, err error) {
	key := prefixVideoURL + objectID

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
			return "", "", err
		}
	}
	if err != nil {
		return "", "", err
	}

	// 乐观锁成功, 结果可信
	if exists == 0 {
		return "", "", ErrorRedisNil
	}
	videoOSS := &videoOSS{}
	err = cmd.Scan(videoOSS)
	if err != nil {
		return "", "", err
	}
	return videoOSS.VideoURL, videoOSS.CoverURL, nil
}
