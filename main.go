package main

import (
	"douyin/conf"
	"douyin/repo/db"
	"douyin/repo/oss"
	"douyin/repo/redis"
	"douyin/router"
	"douyin/utility"

	"context"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func init() {
	// 读取配置并初始化公共日志记录器
	conf.InitConfig()
	utility.InitLogger()

	// 设定Gin模式(GORM的日志记录模式将与Gin模式相同)
	if strings.ToLower(conf.Cfg().Log.Level) != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化存储层
	db.InitMySQL()
	oss.InitOSS()
	redis.InitRedis()
}

func main() {
	// 打印配置内容
	utility.PrintAsJson(conf.Cfg())

	// 自动迁移
	if conf.Cfg().MySQL.AutoMigrate {
		db.MakeMigrate()
	}

	// 处理终止信号(优雅关闭)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,  // CTRL+C
		syscall.SIGTERM, // kill
	)
	defer stop() // 停止处理信号, 而非停机

	// 启动服务
	var err error = nil
	r := router.NewRouter()
	if strings.ToLower(conf.Cfg().System.AutoTLS) != "none" {
		err = autotls.RunWithContext(ctx, r, conf.Cfg().System.AutoTLS)
	} else {
		err = router.RunWithContext(ctx, r, conf.Cfg().System.ListenAddress+":"+conf.Cfg().System.ListenPort)
	}
	if err != nil {
		if err == http.ErrServerClosed {
			utility.Logger().Warnf("main warn: %v", err)
		} else {
			utility.Logger().Fatalf("main ftal: %v", err)
		}
	}
}
