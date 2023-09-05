package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

const prefixRateLimiter = "rate:" // 后接IP

// 处理限流(判断是否放行) (旧)
func CheckRateSimple(ctx context.Context, ip string, limit int, period time.Duration) (ok bool) {
	key := prefixRateLimiter + ip
	new, err := _redis.Incr(ctx, key).Result()
	if err == nil {
		if new == 1 { // Incr前key必定为0或不存在, 流程中无置零操作, 因此此时必为Incr刚刚新建key
			for {
				if _redis.Expire(ctx, key, period).Err() == nil { // 必须成功(可以循环重试, 因为只可能因网络原因失败)
					break
				}
			}
			return true
		} else if new <= int64(limit) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// 处理限流(判断是否放行)
func CheckRate(ctx context.Context, ip string, limit int, period time.Duration) (ok bool) {
	key := prefixRateLimiter + ip
	cfg := redis_rate.Limit{Rate: limit, Period: period, Burst: limit}
	result, err := _limiter.Allow(ctx, key, cfg)
	return err == nil && result.Allowed > 0
}
