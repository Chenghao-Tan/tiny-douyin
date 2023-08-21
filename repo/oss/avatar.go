package oss

import (
	"douyin/emb"
	"douyin/utility"

	"context"
)

// 自定义对象扩展名 需包含"."
const avatarExt = ".webp"

// 获取头像对象在存储桶内所用名称
func getAvatarObjectName(objectID string) (avatarName string) {
	return "user/" + objectID + "_avatar" + avatarExt // 模拟user文件夹
}

// 获取头像对象的短期外链
func GetAvatar(ctx context.Context, objectID string) (avatarURL string, err error) {
	// 头像对象名
	avatarName := getAvatarObjectName(objectID)

	// 获取URL
	avatarURL, err = _oss.getURL(ctx, avatarName)
	if err != nil {
		utility.Logger().Errorf("_oss.getURL (avatar) err: %v", err)
		return "", err
	}

	return avatarURL, nil
}

// 本项目前仅为流式上传默认头像对象
func UploadAvatarStream(ctx context.Context, objectID string) (err error) {
	// 头像对象名
	avatarName := getAvatarObjectName(objectID)

	// 获取默认头像
	avatarStream, err := emb.Emb().Open("assets/defaultAvatar" + avatarExt)
	if err != nil {
		utility.Logger().Errorf("Emb().Open (defaultAvatar) err: %v", err)
		return err
	}
	defer avatarStream.Close() // 不保证自动关闭成功

	avatarStat, err := avatarStream.Stat()
	if err != nil {
		utility.Logger().Errorf("File.Stat (defaultAvatar) err: %v", err)
		return err
	}
	avatarSize := avatarStat.Size()

	err = _oss.uploadStream(ctx, avatarName, avatarStream, avatarSize)
	if err != nil {
		utility.Logger().Errorf("_oss.uploadStream (avatar) err: %v", err)
		return err
	}

	return nil
}
