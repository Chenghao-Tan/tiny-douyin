package db

import (
	"douyin/conf"
	"douyin/repo/internal/db/model"

	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 自定义错误类型
var ErrorRecordExists = errors.New("记录已存在")
var ErrorRecordNotExists = errors.New("记录不存在")

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

	var datetimePrecision int = 0 // 精度为秒
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        256,                // 设定string类型字段的默认长度
		DefaultDatetimePrecision: &datetimePrecision, // 设定datetime精度
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名使用单数形式
		},
		Logger:                                   ormLogger, // 打印日志
		PrepareStmt:                              true,      // 缓存预编译语句
		DisableForeignKeyConstraintWhenMigrating: false,     // 设置是否关闭迁移时自动创建外键约束
	})
	if err != nil {
		panic(err)
	}

	_db = db
}

// 迁移数据表
// 只支持创建表与增加表中没有的字段和索引
// 为了保护数据, 并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate() (err error) {
	DB := _db.WithContext(context.Background())
	return DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{}, &model.Message{})
}
