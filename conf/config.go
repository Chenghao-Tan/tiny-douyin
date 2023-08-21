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

	// 特殊值替换
	if strings.ToLower(_cfg.System.TempDir) == "system" { // 若使用系统默认临时文件夹
		_cfg.System.TempDir = filepath.Join(os.TempDir(), "douyin")
	}
}
