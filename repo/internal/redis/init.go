package redis

import (
	"douyin/conf"
	"douyin/repo/internal/db"

	"crypto/tls"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

// 自定义错误类型
const ErrorRedisNil = redis.Nil              // 创建查询为空的别名
const ErrorRedisTxFailed = redis.TxFailedErr // 创建乐观锁失败的别名
var ErrorRecordExists = db.ErrorRecordExists
var ErrorRecordNotExists = db.ErrorRecordNotExists
var ErrorSelfFollow = db.ErrorSelfFollow

const randomExpirationRatio = 0.1 // 随机延长过期时间的比例(防止缓存雪崩)
const maxRetries = 3              // watch乐观锁等的最大重试次数

var _redis *redis.Client

func InitRedis() {
	redisCfg := conf.Cfg().Redis

	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisCfg.RedisHost, redisCfg.RedisPort),
		DB:   redisCfg.RedisDB,
	}

	if strings.ToLower(redisCfg.Username) != "none" {
		opts.Username = redisCfg.Username
	}

	if strings.ToLower(redisCfg.Password) != "none" {
		opts.Password = redisCfg.Password
	}

	if redisCfg.TLS {
		opts.TLSConfig = &tls.Config{ServerName: redisCfg.RedisHost} // 默认要求TLS最低版本1.2 可通过MinVersion指定
	}

	_redis = redis.NewClient(opts)
	rateLimiter = redis_rate.NewLimiter(_redis)
}

// 随机轻微延长过期时间以防止缓存雪崩
func randomExpiration(expiration time.Duration) (randomized time.Duration) {
	return expiration + time.Duration(rand.Intn(int(float64(expiration)*randomExpirationRatio))).Abs()
}
