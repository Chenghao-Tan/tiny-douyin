package main

import (
	"douyin/conf"
	"douyin/repo/db"
	"douyin/repo/oss"
	"douyin/repo/redis"
	"douyin/router"
	"douyin/utility"

	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func init() {
	conf.InitConfig()
	utility.InitLogger()
	db.InitMySQL()
	oss.InitOSS()
	redis.InitRedis()
	if strings.ToLower(conf.Cfg().Log.Level) != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	utility.PrintAsJson(conf.Cfg())

	if conf.Cfg().MySQL.AutoMigrate {
		db.MakeMigrate()
	}

	r := router.NewRouter()
	if strings.ToLower(conf.Cfg().System.AutoTLS) != "none" {
		utility.Logger().Fatalf("main ftal: %v", autotls.Run(r, conf.Cfg().System.AutoTLS))
	} else {
		utility.Logger().Fatalf("main ftal: %v", r.Run(":"+conf.Cfg().System.HttpPort))
	}
}
