package conf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	System *System `yaml:"system"`
	MySQL  *MySQL  `yaml:"mysql"`
	OSS    *OSS    `yaml:"oss"`
	Redis  *Redis  `yaml:"redis"`
	Cache  *Cache  `yaml:"cache"`
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
	viper.AddConfigPath(workDir)
	viper.AddConfigPath(workDir + "/conf/locale")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&_cfg)
	if err != nil {
		panic(err)
	}

	// 特殊值替换
	if strings.ToLower(_cfg.System.TempDir) == "system" { // 若使用系统默认临时文件夹
		_cfg.System.TempDir = filepath.Join(os.TempDir(), "douyin")
	}
	if strings.HasPrefix(_cfg.MySQL.DbHost, "$") { // 若使用环境变量(即$开头)
		_cfg.MySQL.DbHost = os.Getenv(_cfg.MySQL.DbHost[1:])
	}
	if strings.HasPrefix(_cfg.MySQL.DbPort, "$") { // 若使用环境变量(即$开头)
		_cfg.MySQL.DbPort = os.Getenv(_cfg.MySQL.DbPort[1:])
	}
	if strings.HasPrefix(_cfg.MySQL.Username, "$") { // 若使用环境变量(即$开头)
		_cfg.MySQL.Username = os.Getenv(_cfg.MySQL.Username[1:])
	}
	if strings.HasPrefix(_cfg.MySQL.Password, "$") { // 若使用环境变量(即$开头)
		_cfg.MySQL.Password = os.Getenv(_cfg.MySQL.Password[1:])
	}
	if strings.HasPrefix(_cfg.OSS.OssHost, "$") { // 若使用环境变量(即$开头)
		_cfg.OSS.OssHost = os.Getenv(_cfg.OSS.OssHost[1:])
	}
	if strings.HasPrefix(_cfg.OSS.OssPort, "$") { // 若使用环境变量(即$开头)
		_cfg.OSS.OssPort = os.Getenv(_cfg.OSS.OssPort[1:])
	}
	if strings.HasPrefix(_cfg.OSS.AccessKeyID, "$") { // 若使用环境变量(即$开头)
		_cfg.OSS.AccessKeyID = os.Getenv(_cfg.OSS.AccessKeyID[1:])
	}
	if strings.HasPrefix(_cfg.OSS.SecretAccessKey, "$") { // 若使用环境变量(即$开头)
		_cfg.OSS.SecretAccessKey = os.Getenv(_cfg.OSS.SecretAccessKey[1:])
	}
	if strings.HasPrefix(_cfg.Redis.RedisHost, "$") { // 若使用环境变量(即$开头)
		_cfg.Redis.RedisHost = os.Getenv(_cfg.Redis.RedisHost[1:])
	}
	if strings.HasPrefix(_cfg.Redis.RedisPort, "$") { // 若使用环境变量(即$开头)
		_cfg.Redis.RedisPort = os.Getenv(_cfg.Redis.RedisPort[1:])
	}
	if strings.HasPrefix(_cfg.Redis.Username, "$") { // 若使用环境变量(即$开头)
		_cfg.Redis.Username = os.Getenv(_cfg.Redis.Username[1:])
	}
	if strings.HasPrefix(_cfg.Redis.Password, "$") { // 若使用环境变量(即$开头)
		_cfg.Redis.Password = os.Getenv(_cfg.Redis.Password[1:])
	}
}
