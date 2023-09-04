package oss

import (
	"douyin/emb"
	"douyin/utility"

	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// 自定义错误类型
var ErrorRollbackFailed = errors.New("回滚操作(封面移除)失败")

// 自定义对象扩展名 需包含"."
const videoExt = ".mp4"
const coverExt = ".png"

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
		utility.Logger().Errorf("_oss.getURL (video) err: %v", err)
		return "", "", err
	}
	coverURL, err = _oss.getURL(ctx, coverName)
	if err != nil {
		utility.Logger().Errorf("_oss.getURL (cover) err: %v", err)
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
	coverPath := filepath.Join(ossTempDir, coverName)  // 临时文件位置
	err = utility.GetSnapshot(videoPath, coverPath, 1) // 切取索引为1的帧 防止切取黑屏
	if err != nil {
		utility.Logger().Errorf("GetSnapshot err: %v", err)
		return err
	}
	defer os.Remove(coverPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接覆写

	// 上传
	err = _oss.upload(ctx, coverName, coverPath) // 先传输小文件
	if err != nil {
		utility.Logger().Errorf("_oss.upload (cover) err: %v", err)
		return err
	}
	err = _oss.upload(ctx, videoName, videoPath)
	if err != nil {
		utility.Logger().Errorf("_oss.upload (video) err: %v", err)

		// 视频传输失败时将移除其封面
		utility.Logger().Warnf("_oss.upload (cover) warn: 正在回滚(移除对应封面%v)", coverName)
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
	coverStream, err := emb.Emb().Open("assets/defaultCover" + coverExt)
	if err != nil {
		utility.Logger().Errorf("Emb().Open (defaultCover) err: %v", err)
		return err
	}
	defer coverStream.Close() // 不保证自动关闭成功

	coverStat, err := coverStream.Stat()
	if err != nil {
		utility.Logger().Errorf("File.Stat (defaultCover) err: %v", err)
		return err
	}
	coverSize := coverStat.Size()

	// 上传
	err = _oss.uploadStream(ctx, coverName, coverStream, coverSize) // 先传输小文件
	if err != nil {
		utility.Logger().Errorf("_oss.uploadStream (cover) err: %v", err)
		return err
	}
	err = _oss.uploadStream(ctx, videoName, videoStream, videoSize)
	if err != nil {
		utility.Logger().Errorf("_oss.uploadStream (video) err: %v", err)

		// 视频传输失败时将移除其封面
		utility.Logger().Warnf("_oss.uploadStream (cover) warn: 正在回滚(移除对应封面%v)", coverName)
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

	// 尝试使用云处理切取封面
	err = _oss.setOperation(ctx, OpUpdateCover, videoName, coverName)
	if err != nil {
		if err != ErrorNotSupported { // 若err==ErrorNotSupported则为OSS不支持该云处理操作 忽略并进行后续处理
			// 回报其他功能性错误
			utility.Logger().Errorf("_oss.setOperation (OpUpdateCover) err: %v", err)
			return err
		}
	} else { // 已设定云处理任务
		utility.Logger().Infof("UpdateCover info: %v - 将由云自动处理", coverName)
		return nil // 提前结束
	}

	// 下载视频对象到本地
	videoPath := filepath.Join(ossTempDir, videoName)
	err = _oss.download(ctx, videoName, videoPath)
	if err != nil {
		utility.Logger().Errorf("_oss.download (video) err: %v", err)
		return err
	}
	defer os.Remove(videoPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接覆写

	// 切取封面
	coverPath := filepath.Join(ossTempDir, coverName)  // 临时文件位置
	err = utility.GetSnapshot(videoPath, coverPath, 1) // 切取索引为1的帧 防止切取黑屏
	if err != nil {
		utility.Logger().Errorf("GetSnapshot err: %v", err)
		return err
	}
	defer os.Remove(coverPath) // 不保证自动清理成功 但临时数据在本地 易于检测是否仍存在且可被直接覆写

	// 上传
	err = _oss.upload(ctx, coverName, coverPath)
	if err != nil {
		utility.Logger().Errorf("_oss.upload (cover) err: %v", err)
		return err
	}

	utility.Logger().Infof("UpdateCover info: %v - 操作成功", coverName)
	return nil
}
