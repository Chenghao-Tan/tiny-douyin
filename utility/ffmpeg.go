package utility

import (
	"douyin/conf"

	"bytes"
	"fmt"
	"strings"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// 切取视频第frameNum帧并保存
func GetSnapshot(videoPath string, snapshotPath string, frameNum int) (err error) {
	buf := bytes.NewBuffer(nil)
	task := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf)

	if strings.ToLower(conf.Cfg().System.FFmpeg) != "system" {
		task = task.SetFfmpegPath(conf.Cfg().System.FFmpeg) // 自定义ffmpeg二进制文件位置
	}

	err = task.Run()
	if err != nil {
		Logger().Errorf("ffmpeg err: %v", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		Logger().Errorf("imagingDecode err: %v", err)
		return err
	}

	err = imaging.Save(img, snapshotPath)
	if err != nil {
		Logger().Errorf("imagingSave: %v", err)
		return err
	}

	return nil
}
