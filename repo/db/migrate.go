package db

import (
	"douyin/repo/db/model"
	"douyin/utility"

	"context"
)

// MakeMigrate 迁移数据表，在没有数据表结构变更时候，建议注释不执行
// 只支持创建表、增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate() {
	DB := _db.WithContext(context.Background())
	err := DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{}, &model.Message{})
	if err != nil {
		panic("数据表迁移失败")
	} else {
		utility.Logger().Warnf("AutoMigrate warn: 数据表迁移成功")
	}
}
