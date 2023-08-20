package oss

import (
	"douyin/conf"
	"douyin/utils"

	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 自定义错误类型
var ErrorRollbackFailed = errors.New("回滚操作(封面移除)失败")

// 自定义对象扩展名 需包含"."
const videoExt = ".mp4"
const coverExt = ".png" // 七牛云云切取封面硬限制为.png
const avatarExt = ".webp"
const backgroundImageExt = ".webp"

// 欢迎添加对其他OSS的支持
type OSService interface {
	init()
	upload(ctx context.Context, objectName string, filePath string) (err error)  // 上传对象
	getURL(ctx context.Context, objectName string) (objectURL string, err error) // 获取对象外链
	remove(ctx context.Context, objectName string) (err error)                   // 移除对象

	// 以下为流式上传方案所需
	uploadStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64) (err error) // 流式上传对象
	download(ctx context.Context, objectName string, filePath string) (err error)                        // 下载对象
}

var _oss OSService

func InitOSS() {
	if strings.ToLower(conf.Cfg().OSS.Service) == "minio" {
		_oss = &minIOService{}
	} else if strings.ToLower(conf.Cfg().OSS.Service) == "qiniu" {
		_oss = &qiNiuService{}
	} else {
		panic(errors.New("暂不支持该OSS: " + conf.Cfg().OSS.Service))
	}
	_oss.init()

	// 确保临时路径存在
	err := os.MkdirAll(filepath.Join(conf.Cfg().System.TempDir, "oss", ""), 0755)
	if err != nil {
		panic(err)
	}
}

// 以下为视频与封面相关
// 获取视频对象与封面对象在存储桶内所用名称
func getVideoObjectName(objectID string) (videoName string, coverName string) {
	return "video/" + objectID + "_video" + videoExt, "video/" + objectID + "_cover" + coverExt // 模拟video文件夹
}

// 获取视频对象与封面对象的短期外链
func GetVideo(ctx context.Context, objectID string) (videoURL string, coverURL string, err error) {
	// 视频对象与封面对象名
	videoName, coverName := getVideoObjectName(objectID)

	// 获取URL
	videoURL, err = _oss.getURL(ctx, videoName)
	if err != nil {
		utils.Logger().Errorf("_oss.getURL (video) err: %v", err)
		return "", "", err
	}
	coverURL, err = _oss.getURL(ctx, coverName)
	if err != nil {
		utils.Logger().Errorf("_oss.getURL (cover) err: %v", err)
		return videoURL, "", err
	}

	return videoURL, coverURL, nil
}

// 文件上传方案已弃用 请使用流式上传方案
// 上传视频对象 自动切取并上传封面对象
func UploadVideo(ctx context.Context, objectID string, videoPath string) (err error) {
	// 视频对象与封面对象名
	videoName, coverName := getVideoObjectName(objectID)

	// 切取封面
	coverPath := filepath.Join(conf.Cfg().System.TempDir, "oss", coverName) // 临时文件位置
	err = utils.GetSnapshot(videoPath, coverPath, 1)                        // 切取索引为1的帧 防止切取黑屏
	if err != nil {
		utils.Logger().Errorf("GetSnapshot err: %v", err)
		return err
	}
	defer os.Remove(coverPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接覆写

	// 上传
	err = _oss.upload(ctx, coverName, coverPath) // 先传输小文件
	if err != nil {
		utils.Logger().Errorf("_oss.upload (cover) err: %v", err)
		return err
	}
	err = _oss.upload(ctx, videoName, videoPath)
	if err != nil {
		utils.Logger().Errorf("_oss.upload (video) err: %v", err)

		// 视频传输失败时将移除其封面
		utils.Logger().Warnf("_oss.upload (cover) warn: 正在回滚(移除对应封面%v)", coverName)
		err2 := _oss.remove(ctx, coverName)
		if err2 != nil {
			return ErrorRollbackFailed
		} else {
			return err
		}
	}

	return nil
}

// 以下为流式上传方案所需
// 流式上传视频对象 自动上传默认封面对象
func UploadVideoStream(ctx context.Context, objectID string, videoStream io.Reader, videoSize int64) (err error) {
	// 视频对象与封面对象名
	videoName, coverName := getVideoObjectName(objectID)

	// 获取默认封面
	coverStream, err := conf.Emb().Open("assets/defaultCover" + coverExt)
	if err != nil {
		utils.Logger().Errorf("Emb().Open (defaultCover) err: %v", err)
		return err
	}
	defer coverStream.Close() // 不保证自动关闭成功

	coverStat, err := coverStream.Stat()
	if err != nil {
		utils.Logger().Errorf("File.Stat (defaultCover) err: %v", err)
		return err
	}
	coverSize := coverStat.Size()

	// 上传
	err = _oss.uploadStream(ctx, coverName, coverStream, coverSize) // 先传输小文件
	if err != nil {
		utils.Logger().Errorf("_oss.uploadStream (cover) err: %v", err)
		return err
	}
	err = _oss.uploadStream(ctx, videoName, videoStream, videoSize)
	if err != nil {
		utils.Logger().Errorf("_oss.uploadStream (video) err: %v", err)

		// 视频传输失败时将移除其封面
		utils.Logger().Warnf("_oss.uploadStream (cover) warn: 正在回滚(移除对应封面%v)", coverName)
		err2 := _oss.remove(ctx, coverName)
		if err2 != nil {
			return ErrorRollbackFailed
		} else {
			return err
		}
	}

	return nil
}

// 更新封面
func UpdateCover(ctx context.Context, objectID string) (err error) {
	// 视频对象与封面对象名
	videoName, coverName := getVideoObjectName(objectID)

	// 七牛云等带有云切取的OSS特殊处理
	if strings.ToLower(conf.Cfg().OSS.Service) == "qiniu" {
		utils.Logger().Infof("UpdateCover info: %v - 将由云自动处理", coverName)
		return nil
	}

	// 下载视频对象到本地
	videoPath := filepath.Join(conf.Cfg().System.TempDir, "oss", videoName)
	err = _oss.download(ctx, videoName, videoPath)
	if err != nil {
		utils.Logger().Errorf("_oss.download (video) err: %v", err)
		return err
	}
	defer os.Remove(videoPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接覆写

	// 切取封面
	coverPath := filepath.Join(conf.Cfg().System.TempDir, "oss", coverName) // 临时文件位置
	err = utils.GetSnapshot(videoPath, coverPath, 1)                        // 切取索引为1的帧 防止切取黑屏
	if err != nil {
		utils.Logger().Errorf("GetSnapshot err: %v", err)
		return err
	}
	defer os.Remove(coverPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接覆写

	// 上传
	err = _oss.upload(ctx, coverName, coverPath)
	if err != nil {
		utils.Logger().Errorf("_oss.upload (cover) err: %v", err)
		return err
	}

	utils.Logger().Infof("UpdateCover info: %v - 操作成功", coverName)
	return nil
}

// 以下为头像与个人页背景图相关 流式上传以外的方案不受支持
// 获取头像对象在存储桶内所用名称
func getAvatarObjectName(objectID string) (avatarName string) {
	return "user/" + objectID + "_avatar" + avatarExt // 模拟user文件夹
}

// 获取个人页背景图对象在存储桶内所用名称
func getBackgroundImageObjectName(objectID string) (backgroundImageName string) {
	return "user/" + objectID + "_backgroundImage" + backgroundImageExt // 模拟user文件夹
}

// 获取头像对象的短期外链
func GetAvatar(ctx context.Context, objectID string) (avatarURL string, err error) {
	// 头像对象名
	avatarName := getAvatarObjectName(objectID)

	// 获取URL
	avatarURL, err = _oss.getURL(ctx, avatarName)
	if err != nil {
		utils.Logger().Errorf("_oss.getURL (avatar) err: %v", err)
		return "", err
	}

	return avatarURL, nil
}

// 本项目前仅为流式上传默认头像对象
func UploadAvatarStream(ctx context.Context, objectID string) (err error) {
	// 头像对象名
	avatarName := getAvatarObjectName(objectID)

	// 获取默认头像
	avatarStream, err := conf.Emb().Open("assets/defaultAvatar" + avatarExt)
	if err != nil {
		utils.Logger().Errorf("Emb().Open (defaultAvatar) err: %v", err)
		return err
	}
	defer avatarStream.Close() // 不保证自动关闭成功

	avatarStat, err := avatarStream.Stat()
	if err != nil {
		utils.Logger().Errorf("File.Stat (defaultAvatar) err: %v", err)
		return err
	}
	avatarSize := avatarStat.Size()

	err = _oss.uploadStream(ctx, avatarName, avatarStream, avatarSize)
	if err != nil {
		utils.Logger().Errorf("_oss.uploadStream (avatar) err: %v", err)
		return err
	}

	return nil
}

// 获取个人页背景图对象的短期外链
func GetBackgroundImage(ctx context.Context, objectID string) (backgroundImageURL string, err error) {
	// 个人页背景图对象名
	backgroundImageName := getBackgroundImageObjectName(objectID)

	// 获取URL
	backgroundImageURL, err = _oss.getURL(ctx, backgroundImageName)
	if err != nil {
		utils.Logger().Errorf("_oss.getURL (backgroundImage) err: %v", err)
		return "", err
	}

	return backgroundImageURL, nil
}

// 本项目前仅为流式上传默认个人页背景图对象
func UploadBackgroundImageStream(ctx context.Context, objectID string) (err error) {
	// 头像对象名
	backgroundImageName := getBackgroundImageObjectName(objectID)

	// 获取默认头像
	backgroundImageStream, err := conf.Emb().Open("assets/defaultBackgroundImage" + backgroundImageExt)
	if err != nil {
		utils.Logger().Errorf("Emb().Open (defaultBackgroundImage) err: %v", err)
		return err
	}
	defer backgroundImageStream.Close() // 不保证自动关闭成功

	backgroundImageStat, err := backgroundImageStream.Stat()
	if err != nil {
		utils.Logger().Errorf("File.Stat (defaultBackgroundImage) err: %v", err)
		return err
	}
	backgroundImageSize := backgroundImageStat.Size()

	err = _oss.uploadStream(ctx, backgroundImageName, backgroundImageStream, backgroundImageSize)
	if err != nil {
		utils.Logger().Errorf("_oss.uploadStream (backgroundImage) err: %v", err)
		return err
	}

	return nil
}
