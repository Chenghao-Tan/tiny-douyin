package repo

import (
	"douyin/conf"
	"douyin/repo/internal/db"
	"douyin/repo/internal/oss"
	"douyin/repo/internal/redis"
	"douyin/utility"

	"errors"
	"time"
)

// 自定义错误类型
var ErrorEmptyObject = errors.New("对象不存在或尚不存在")

var syncInterval time.Duration
var maxRWTime time.Duration
var cacheExpiration time.Duration
var emptyExpiration time.Duration
var distrustProbability float32
var urlExpiration time.Duration

func Init() {
	cacheCfg := conf.Cfg().Cache
	syncInterval = time.Second * time.Duration(cacheCfg.SyncInterval).Abs()
	maxRWTime = time.Millisecond * time.Duration(cacheCfg.MaxRWTime).Abs()
	cacheExpiration = time.Second * time.Duration(cacheCfg.CacheExpiration).Abs()
	emptyExpiration = time.Second * time.Duration(cacheCfg.EmptyExpiration).Abs()
	distrustProbability = cacheCfg.DistrustProbability
	urlExpiration = time.Hour*time.Duration(conf.Cfg().OSS.Expiry).Abs() - time.Minute

	// 初始化存储层
	db.InitMySQL()
	oss.InitOSS()
	redis.InitRedis()

	// 自动迁移
	if conf.Cfg().MySQL.AutoMigrate {
		err := db.MakeMigrate()
		if err != nil {
			panic(err)
		} else {
			utility.Logger().Warnf("repo.Init warn: 数据表迁移成功")
		}
	}

	// 初始化同步系统
	syncQueue.Init()
	syncCron.AddFunc("@every "+syncInterval.String(), syncTask)
	syncCron.Start()
}

func Stop() {
	syncCron.Stop()
	utility.Logger().Warnf("repo.Stop warn: 已停止启动新任务, 正在等待现有任务结束...")
	time.Sleep(maxRWTime) // 等待同步任务(如有)彻底结束
}
