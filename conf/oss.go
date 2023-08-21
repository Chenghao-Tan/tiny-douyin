package conf

type OSS struct {
	Service         string `yaml:"service"`
	OssHost         string `yaml:"ossHost"`
	OssPort         string `yaml:"ossPort"`
	OssRegion       string `yaml:"ossRegion"`
	BucketName      string `yaml:"bucketName"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	TLS             bool   `yaml:"tls"`
	Expiry          int    `yaml:"expiry"`
}
