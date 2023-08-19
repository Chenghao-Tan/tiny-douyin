package main

import (
	"douyin/conf"
	"douyin/repository/dao"
	"douyin/routes"
	"douyin/utils"
	"douyin/utils/oss"

	"strings"

	"github.com/gin-gonic/autotls"
)

func init() {
	conf.InitConfig()
	utils.InitLogger()
	dao.InitMySQL()
	oss.InitOSS()
}

func main() {
	utils.PrintAsJson(conf.Cfg())

	// 初次使用或数据表结构变更时取消以下行的注释以迁移数据表
	// dao.MakeMigrate()

	r := routes.NewRouter()
	if strings.ToLower(conf.Cfg().System.AutoTLS) != "none" {
		utils.Logger().Fatalf("main ftal: %v", autotls.Run(r, conf.Cfg().System.AutoTLS))
	} else {
		utils.Logger().Fatalf("main ftal: %v", r.Run(":"+conf.Cfg().System.HttpPort))
	}
}
