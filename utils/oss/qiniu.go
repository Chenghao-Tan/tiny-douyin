package oss

import (
	"douyin/conf"

	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type qiNiuService struct {
	mac        *auth.Credentials
	cfg        *storage.Config
	bucketName string
	domain     string
	expires    time.Duration
}

func isVideo(objectName string) (is bool) {
	videoExts := []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".webm"} // 常见视频文件扩展名列表

	dotIndex := strings.LastIndex(objectName, ".")
	if dotIndex == -1 || dotIndex == len(objectName)-1 {
		return false // 没有有效的扩展名
	}

	ext := strings.ToLower(objectName[dotIndex:])
	for _, videoExt := range videoExts {
		if ext == videoExt {
			return true
		}
	}

	return false
}

func (q *qiNiuService) init() {
	accessKey := conf.Cfg().OSS.AccessKeyID
	secretKey := conf.Cfg().OSS.SecretAccessKey
	q.mac = qbox.NewMac(accessKey, secretKey)

	q.cfg = &storage.Config{
		UseHTTPS:      conf.Cfg().OSS.TLS, // 是否使用TLS
		UseCdnDomains: true,               // 是否使用CDN上传加速
	}

	// 空间对应的机房
	switch strings.ToLower(conf.Cfg().OSS.OssRegion) {
	case "huadong":
		q.cfg.Region = &storage.ZoneHuadong
	case "huadongzhejiang2":
		q.cfg.Region = &storage.ZoneHuadongZheJiang2
	case "huabei":
		q.cfg.Region = &storage.ZoneHuabei
	case "huanan":
		q.cfg.Region = &storage.ZoneHuanan
	case "beimei":
		q.cfg.Region = &storage.ZoneBeimei
	case "xinjiapo":
		q.cfg.Region = &storage.ZoneXinjiapo
	case "shouer":
		q.cfg.Region = &storage.ZoneShouEr1
	}

	q.bucketName = conf.Cfg().OSS.BucketName
	q.domain = conf.Cfg().OSS.OssHost
	q.expires = time.Hour * time.Duration(conf.Cfg().OSS.Expiry).Abs()
}

// 获取自动云切取PutPolicy或普通PutPolicy
func (q *qiNiuService) getPutPolicy(objectName string, snapshot bool) (putPolicy *storage.PutPolicy) {
	if snapshot {
		// 构建封面云切取操作
		coverName := strings.Split(objectName, ".")[0] + ".png"                                                                                  // 硬限制为png格式
		vframePngFop := "vframe/png/offset/1|saveas/" + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", q.bucketName, coverName))) // 切取索引为1的帧 防止切取黑屏
		persistentOps := strings.Join([]string{vframePngFop}, ";")                                                                               // 仅使用云切取指令
		pipeline := ""                                                                                                                           // 使用公有队列

		return &storage.PutPolicy{
			Scope:              q.bucketName + ":" + objectName, // 设定为允许覆盖
			PersistentOps:      persistentOps,
			PersistentPipeline: pipeline,
		}
	} else {
		return &storage.PutPolicy{
			Scope: q.bucketName + ":" + objectName, // 设定为允许覆盖
		}
	}
}

// 若上传为常见格式视频则将自动云切取同名.png封面!!!
func (q *qiNiuService) upload(ctx context.Context, objectName string, filePath string) (err error) {
	putPolicy := q.getPutPolicy(objectName, isVideo(objectName))
	upToken := putPolicy.UploadToken(q.mac)        // token有效期默认1小时
	formUploader := storage.NewFormUploader(q.cfg) // 构建表单上传对象
	return formUploader.PutFile(ctx, &storage.PutRet{}, upToken, objectName, filePath, &storage.PutExtra{})
}

// SDK限制 context不可用
func (q *qiNiuService) getURL(ctx context.Context, objectName string) (objectURL string, err error) {
	deadline := time.Now().Add(q.expires).Unix() // 此时间后URL失效
	url := storage.MakePrivateURL(q.mac, q.domain, objectName, deadline)

	// 添加协议类型
	if q.cfg.UseHTTPS {
		url = "https://" + url
	} else {
		url = "http://" + url
	}

	return url, nil // SDK限制 始终回报成功
}

// SDK限制 context不可用
func (q *qiNiuService) remove(ctx context.Context, objectName string) (err error) {
	bucketManager := storage.NewBucketManager(q.mac, q.cfg)
	return bucketManager.Delete(q.bucketName, objectName)
}

// 若上传为常见格式视频则将自动云切取同名.png封面!!!
func (q *qiNiuService) uploadStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64) (err error) {
	putPolicy := q.getPutPolicy(objectName, isVideo(objectName))
	upToken := putPolicy.UploadToken(q.mac)        // token有效期默认1小时
	formUploader := storage.NewFormUploader(q.cfg) // 构建表单上传对象
	return formUploader.Put(ctx, &storage.PutRet{}, upToken, objectName, reader, objectSize, &storage.PutExtra{})
}

// 仅请求过程context可用
func (q *qiNiuService) download(ctx context.Context, objectName string, filePath string) (err error) {
	deadline := time.Now().Add(time.Hour).Unix() // 需在1小时内下载完毕
	url := storage.MakePrivateURL(q.mac, q.domain, objectName, deadline)

	// 发送GET请求以下载
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close() // 不保证自动关闭成功

	if response.StatusCode != http.StatusOK { // 若响应非请求成功
		return errors.New("请求失败")
	}

	// 创建输出文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close() // 不保证自动关闭成功

	// 流式写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
