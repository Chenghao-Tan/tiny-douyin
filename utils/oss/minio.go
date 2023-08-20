package oss

import (
	"douyin/conf"

	"context"
	"io"
	"net"
	"net/url"
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
	endpoint := net.JoinHostPort(conf.Cfg().OSS.OssHost, conf.Cfg().OSS.OssPort)
	accessKeyID := conf.Cfg().OSS.AccessKeyID
	secretAccessKey := conf.Cfg().OSS.SecretAccessKey

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: conf.Cfg().OSS.TLS,       // 设置是否使用TLS访问对象存储
		Region: conf.Cfg().OSS.OssRegion, // 设置区域(留空则默认为us-east-1)
	})
	if err != nil {
		panic(err)
	}

	m.client = client
	m.bucketName = conf.Cfg().OSS.BucketName
	m.expires = time.Hour * time.Duration(conf.Cfg().OSS.Expiry).Abs()
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
