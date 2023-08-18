package main

import (
	"douyin/conf"
	"douyin/repository/dao"
	"douyin/routes"
	"douyin/utils"
	"douyin/utils/oss"
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
	_ = r.Run(":" + conf.Cfg().System.HttpPort)
}
