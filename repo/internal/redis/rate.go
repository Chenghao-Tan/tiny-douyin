package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const prefixRateLimiter = "rate:" // 后接IP

// 处理限流(判断是否放行)
func CheckRate(ctx context.Context, ip string, limit int, period time.Duration) (ok bool) {
	key := prefixRateLimiter + ip

	cmds, err := _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		pipe.Get(ctx, key)
		pipe.Incr(ctx, key)
		return nil
	})
	if err != nil {
		return false
	}

	current, err := cmds[0].(*redis.StringCmd).Int() // 读取pipe.Get的结果
	if err == nil {
		if current < limit {
			return true
		} else {
			return false
		}
	} else if err == ErrorRedisNil {
		_ = _redis.Expire(ctx, key, period).Err() // 因为pipe.Incr所以key大概率存在
		return true
	} else {
		return false
	}
}
