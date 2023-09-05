package repo

import (
	"douyin/repo/internal/oss"
	"douyin/repo/internal/redis"

	"context"
	"io"
	"time"
)

// 获取视频对象与封面对象的短期外链
func GetVideo(ctx context.Context, objectID string) (videoURL string, coverURL string, err error) {
	videoURL, coverURL, err = redis.GetVideoURL(ctx, objectID)
	if err == nil { // 命中缓存
		if videoURL == "" { // 命中空对象
			time.Sleep(maxRWTime)
			videoURL, coverURL, err = redis.GetVideoURL(ctx, objectID) // 重试
		} else {
			return videoURL, coverURL, nil
		}
	}
	if err == nil { // 命中缓存
		if videoURL == "" { // 命中空对象
			return "", "", ErrorEmptyObject
		} else {
			return videoURL, coverURL, nil
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetVideoURL(ctx, objectID, "", "", emptyExpiration) // 防止缓存穿透与缓存击穿
		newVideoURL, newCoverURL, err := oss.GetVideo(ctx, objectID)
		if err == nil {
			_ = redis.SetVideoURL(ctx, objectID, newVideoURL, newCoverURL, urlExpiration)
			return newVideoURL, newCoverURL, nil
		} else {
			return "", "", err
		}
	} else {
		return "", "", err
	}
}

// 流式上传视频对象 自动上传默认封面对象
func UploadVideoStream(ctx context.Context, objectID string, videoStream io.Reader, videoSize int64) (err error) {
	return oss.UploadVideoStream(ctx, objectID, videoStream, videoSize)
}

// 更新封面
func UpdateCover(ctx context.Context, objectID string) (err error) {
	return oss.UpdateCover(ctx, objectID)
}

// 获取头像对象的短期外链
func GetAvatar(ctx context.Context, objectID string) (avatarURL string, err error) {
	avatarURL, err = redis.GetUserAvatarURL(ctx, objectID)
	if err == nil { // 命中缓存
		if avatarURL == "" { // 命中空对象
			time.Sleep(maxRWTime)
			avatarURL, err = redis.GetUserAvatarURL(ctx, objectID) // 重试
		} else {
			return avatarURL, nil
		}
	}
	if err == nil { // 命中缓存
		if avatarURL == "" { // 命中空对象
			return "", ErrorEmptyObject
		} else {
			return avatarURL, nil
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserAvatarURL(ctx, objectID, "", emptyExpiration) // 防止缓存穿透与缓存击穿
		newAvatarURL, err := oss.GetAvatar(ctx, objectID)
		if err == nil {
			_ = redis.SetUserAvatarURL(ctx, objectID, newAvatarURL, urlExpiration)
			return newAvatarURL, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

// 本项目前仅为流式上传默认头像对象
func UploadAvatarStream(ctx context.Context, objectID string) (err error) {
	return oss.UploadAvatarStream(ctx, objectID)
}

// 获取个人页背景图对象的短期外链
func GetBackgroundImage(ctx context.Context, objectID string) (backgroundImageURL string, err error) {
	backgroundImageURL, err = redis.GetUserBackgroundImageURL(ctx, objectID)
	if err == nil { // 命中缓存
		if backgroundImageURL == "" { // 命中空对象
			time.Sleep(maxRWTime)
			backgroundImageURL, err = redis.GetUserBackgroundImageURL(ctx, objectID) // 重试
		} else {
			return backgroundImageURL, nil
		}
	}
	if err == nil { // 命中缓存
		if backgroundImageURL == "" { // 命中空对象
			return "", ErrorEmptyObject
		} else {
			return backgroundImageURL, nil
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserBackgroundImageURL(ctx, objectID, "", emptyExpiration) // 防止缓存穿透与缓存击穿
		newBackgroundImageURL, err := oss.GetBackgroundImage(ctx, objectID)
		if err == nil {
			_ = redis.SetUserBackgroundImageURL(ctx, objectID, newBackgroundImageURL, urlExpiration)
			return newBackgroundImageURL, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

// 本项目前仅为流式上传默认个人页背景图对象
func UploadBackgroundImageStream(ctx context.Context, objectID string) (err error) {
	return oss.UploadBackgroundImageStream(ctx, objectID)
}
