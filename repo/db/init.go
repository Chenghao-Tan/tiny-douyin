package db

import (
	"douyin/conf"

	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

func InitMySQL() {
	mysqlCfg := conf.Cfg().MySQL

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s&charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=True", // 禁止使用BIG5/CP932/GB2312/GBK/SJIS
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.DbHost,
		mysqlCfg.DbPort,
		mysqlCfg.DbName,
		mysqlCfg.TLS,
	)

	var ormLogger = logger.Default
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名使用单数形式
		},
		Logger:      ormLogger, // 打印日志
		PrepareStmt: true,      // 缓存预编译语句
	})
	if err != nil {
		panic(err)
	}

	_db = db
}
