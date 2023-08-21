package oss

import (
	"douyin/emb"
	"douyin/utility"

	"context"
)

// 自定义对象扩展名 需包含"."
const backgroundImageExt = ".webp"

// 获取个人页背景图对象在存储桶内所用名称
func getBackgroundImageObjectName(objectID string) (backgroundImageName string) {
	return "user/" + objectID + "_backgroundImage" + backgroundImageExt // 模拟user文件夹
}

// 获取个人页背景图对象的短期外链
func GetBackgroundImage(ctx context.Context, objectID string) (backgroundImageURL string, err error) {
	// 个人页背景图对象名
	backgroundImageName := getBackgroundImageObjectName(objectID)

	// 获取URL
	backgroundImageURL, err = _oss.getURL(ctx, backgroundImageName)
	if err != nil {
		utility.Logger().Errorf("_oss.getURL (backgroundImage) err: %v", err)
		return "", err
	}

	return backgroundImageURL, nil
}

// 本项目前仅为流式上传默认个人页背景图对象
func UploadBackgroundImageStream(ctx context.Context, objectID string) (err error) {
	// 头像对象名
	backgroundImageName := getBackgroundImageObjectName(objectID)

	// 获取默认头像
	backgroundImageStream, err := emb.Emb().Open("assets/defaultBackgroundImage" + backgroundImageExt)
	if err != nil {
		utility.Logger().Errorf("Emb().Open (defaultBackgroundImage) err: %v", err)
		return err
	}
	defer backgroundImageStream.Close() // 不保证自动关闭成功

	backgroundImageStat, err := backgroundImageStream.Stat()
	if err != nil {
		utility.Logger().Errorf("File.Stat (defaultBackgroundImage) err: %v", err)
		return err
	}
	backgroundImageSize := backgroundImageStat.Size()

	err = _oss.uploadStream(ctx, backgroundImageName, backgroundImageStream, backgroundImageSize)
	if err != nil {
		utility.Logger().Errorf("_oss.uploadStream (backgroundImage) err: %v", err)
		return err
	}

	return nil
}
