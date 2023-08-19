package conf

import (
	"os"

	"github.com/spf13/viper"
)

type System struct {
	HttpPort string `yaml:"httpPort"`
	AutoTLS  string `yaml:"autoTLS"`
	TempDir  string `yaml:"tempDir"`
}

type MySQL struct {
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	TLS      string `yaml:"tls"`
	Charset  string `yaml:"charset"`
}

type OSS struct {
	Service         string `yaml:"service"`
	OssHost         string `yaml:"ossHost"`
	OssPort         string `yaml:"ossPort"`
	BucketName      string `yaml:"bucketName"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	TLS             bool   `yaml:"tls"`
	Expiry          int    `yaml:"expiry"`
}

type Log struct {
	Path       string `yaml:"path"`       // 输出路径
	Level      string `yaml:"level"`      // 输出级别
	Prefix     string `yaml:"prefix"`     // 日志前缀
	ShowLine   bool   `yaml:"showLine"`   // 显示行号
	MaxSize    int    `yaml:"maxSize"`    // 单个日志文件最大大小
	MaxBackups int    `yaml:"maxBackups"` // 最多保留数量
	MaxAge     int    `yaml:"maxAge"`     // 最多保留天数
	Compress   bool   `yaml:"compress"`   // 是否gzip压缩
}

type Config struct {
	System *System `yaml:"system"`
	MySQL  *MySQL  `yaml:"mysql"`
	OSS    *OSS    `yaml:"oss"`
	Log    *Log    `yaml:"log"`
}

var _cfg *Config

func Cfg() *Config {
	return _cfg
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/conf/locale")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&_cfg)
	if err != nil {
		panic(err)
	}
}
