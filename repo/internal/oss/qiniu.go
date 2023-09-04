package oss

import (
	"douyin/conf"

	"context"
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

func (q *qiNiuService) init() {
	ossCfg := conf.Cfg().OSS

	accessKey := ossCfg.AccessKeyID
	secretKey := ossCfg.SecretAccessKey
	q.mac = qbox.NewMac(accessKey, secretKey)

	q.cfg = &storage.Config{
		UseHTTPS:      ossCfg.TLS, // 是否使用TLS
		UseCdnDomains: true,       // 是否使用CDN上传加速
	}

	// 空间对应的机房
	switch strings.ToLower(ossCfg.OssRegion) {
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

	q.bucketName = ossCfg.BucketName
	q.domain = ossCfg.OssHost
	q.expires = time.Hour * time.Duration(ossCfg.Expiry).Abs()
}

func (q *qiNiuService) upload(ctx context.Context, objectName string, filePath string) (err error) {
	putPolicy := &storage.PutPolicy{
		Scope: q.bucketName + ":" + objectName, // 设定为允许覆盖
	}
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

func (q *qiNiuService) uploadStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64) (err error) {
	putPolicy := &storage.PutPolicy{
		Scope: q.bucketName + ":" + objectName, // 设定为允许覆盖
	}
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

func (q *qiNiuService) setOperation(ctx context.Context, operation int, from string, to string) (err error) {
	if operation == OpUpdateCover { // 云切取封面操作
		// 分析目标格式
		supportedExts := []string{".png", ".jpg", ".jpeg"} // 暂只支持这些格式

		dotIndex := strings.LastIndex(to, ".")
		if dotIndex == -1 || dotIndex == len(to)-1 {
			return errors.New("没有有效的扩展名")
		}

		coverExt := strings.ToLower(to[dotIndex:]) // 全小写

		// 云切取为自动探测到的格式
		isSupported := false
		for _, supportedExt := range supportedExts {
			if coverExt == supportedExt {
				isSupported = true
			}
		}
		if isSupported {
			// 设定云切取任务
			saveEntry := storage.EncodedEntry(q.bucketName, to)
			vframeFop := fmt.Sprintf("vframe/%s/offset/1|saveas/%s", coverExt[1:], saveEntry) // 切取索引为1的帧 防止切取黑屏
			persistentOps := strings.Join([]string{vframeFop}, ";")                           // 仅使用云切取指令
			pipeline := ""                                                                    // 使用公有队列
			operationManager := storage.NewOperationManager(q.mac, q.cfg)
			_, err := operationManager.Pfop(q.bucketName, from, persistentOps, pipeline, "", true)
			if err != nil {
				return err
			}
		} else { // 若不支持目标格式
			return ErrorNotSupported // 返回指定错误
		}
	} else { // 不支持其他云处理操作
		return ErrorNotSupported // 返回指定错误
	}

	return nil
}
