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
	// "github.com/gin-gonic/gin"
)

func init() {
	conf.InitConfig()
	utility.InitLogger()
	db.InitMySQL()
	oss.InitOSS()
	redis.InitRedis()
}

func main() {
	utility.PrintAsJson(conf.Cfg())

	// 生产环境时取消以下行的注释
	// gin.SetMode(gin.ReleaseMode)

	// 初次使用或数据表结构变更时取消以下行的注释以迁移数据表
	// db.MakeMigrate()

	r := router.NewRouter()
	if strings.ToLower(conf.Cfg().System.AutoTLS) != "none" {
		utility.Logger().Fatalf("main ftal: %v", autotls.Run(r, conf.Cfg().System.AutoTLS))
	} else {
		utility.Logger().Fatalf("main ftal: %v", r.Run(":"+conf.Cfg().System.HttpPort))
	}
}
