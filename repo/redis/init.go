package redis

import (
	"douyin/conf"

	"crypto/tls"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

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
}
