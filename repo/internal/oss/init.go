package oss

import (
	"douyin/conf"

	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 自定义错误类型
var ErrorNotSupported = errors.New("OSS不支持该操作")

// 自定义云处理操作类型
var OpUpdateCover = 1 // 云切取封面

// 欢迎添加对其他OSS的支持
type OSService interface {
	init()
	upload(ctx context.Context, objectName string, filePath string) (err error)  // 上传对象
	getURL(ctx context.Context, objectName string) (objectURL string, err error) // 获取对象外链
	remove(ctx context.Context, objectName string) (err error)                   // 移除对象

	// 以下为流式上传方案所需
	uploadStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64) (err error) // 流式上传对象
	download(ctx context.Context, objectName string, filePath string) (err error)                        // 下载对象

	// 以下为设定云端处理任务 不兼容OSS应直接返回ErrorNotSupported
	setOperation(ctx context.Context, operation int, from string, to string) (err error) // 设定云端处理任务 operation为操作类型 from和to分别为源对象名和目标对象名
}

var _oss OSService

var ossTempDir string // OSS专属临时路径

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
	ossTempDir = filepath.Join(conf.Cfg().System.TempDir, "oss", "")
	err := os.MkdirAll(ossTempDir, 0755)
	if err != nil {
		panic(err)
	}
}
