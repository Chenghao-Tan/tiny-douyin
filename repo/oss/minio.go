package oss

import (
	"douyin/conf"

	"context"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minIOService struct {
	client     *minio.Client
	bucketName string
	expires    time.Duration
}

func (m *minIOService) init() {
	ossCfg := conf.Cfg().OSS

	endpoint := net.JoinHostPort(ossCfg.OssHost, ossCfg.OssPort)
	accessKeyID := ossCfg.AccessKeyID
	secretAccessKey := ossCfg.SecretAccessKey
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: ossCfg.TLS, // 设置是否使用TLS访问对象存储
	}
	if strings.ToLower(ossCfg.OssRegion) != "default" { // 设置区域
		opts.Region = ossCfg.OssRegion
	}

	client, err := minio.New(endpoint, opts)
	if err != nil {
		panic(err)
	}

	m.client = client
	m.bucketName = ossCfg.BucketName
	m.expires = time.Hour * time.Duration(ossCfg.Expiry).Abs()
}

func (m *minIOService) upload(ctx context.Context, objectName string, filePath string) (err error) {
	opts := minio.PutObjectOptions{} // 可选元数据
	_, err = m.client.FPutObject(ctx, m.bucketName, objectName, filePath, opts)
	return err
}

func (m *minIOService) getURL(ctx context.Context, objectName string) (objectURL string, err error) {
	reqParams := make(url.Values) // 可选响应头
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, m.expires, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *minIOService) remove(ctx context.Context, objectName string) (err error) {
	opts := minio.RemoveObjectOptions{} // 可选选项
	return m.client.RemoveObject(ctx, m.bucketName, objectName, opts)
}

// 若对象大小未知则objectSize可以为-1 但将会占用大量内存!!!
func (m *minIOService) uploadStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64) (err error) {
	opts := minio.PutObjectOptions{} // 可选元数据
	_, err = m.client.PutObject(ctx, m.bucketName, objectName, reader, objectSize, opts)
	return err
}

func (m *minIOService) download(ctx context.Context, objectName string, filePath string) (err error) {
	opts := minio.GetObjectOptions{} // 可选元数据
	return m.client.FGetObject(ctx, m.bucketName, objectName, filePath, opts)
}

func (m *minIOService) setOperation(ctx context.Context, operation int, from string, to string) (err error) {
	return ErrorNotSupported // 返回指定错误 不支持任何云处理操作
}
