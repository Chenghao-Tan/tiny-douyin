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
	videoURL string `redis:"video"`
	coverURL string `redis:"cover"`
}

// 设置视频对象及封面对象外链
func SetVideoURL(ctx context.Context, objectID string, videoURL string, coverURL string, urlExpiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		key := prefixVideoURL + objectID
		value := &videoOSS{videoURL: videoURL, coverURL: coverURL}

		pipe.HSet(ctx, key, value)
		pipe.Expire(ctx, key, urlExpiration)

		return nil
	})
	return err
}

// 读取视频对象及封面对象外链
func GetVideoURL(ctx context.Context, objectID string) (videoURL string, coverURL string, err error) {
	key := prefixVideoURL + objectID
	videoOSS := &videoOSS{}
	err = _redis.HGetAll(ctx, key).Scan(videoOSS)
	if err != nil {
		return "", "", err
	}
	return videoOSS.videoURL, videoOSS.coverURL, nil
}
