package midware

import (
	"douyin/repo"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// gin中间件
// 限流(capacity: 起始/最大请求处理量, recover: 每秒恢复量)
func MiddlewareRateLimit(capacity int64, recover int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(time.Second, capacity, recover)
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			utility.Logger().Warnf("MiddlewareRateLimit warn: 达到处理量上限")
			ctx.JSON(http.StatusTooManyRequests, &response.Status{
				Status_Code: -1,
				Status_Msg:  "请求过频",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// gin中间件
// 使用Redis基于IP地址限流
func MiddlewareRateLimitWithRedis(limit int, period time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if !repo.CheckRate(context.TODO(), ip, limit, period) {
			utility.Logger().Warnf("MiddlewareRateLimitWithRedis warn: %v达到请求量上限", ip)
			ctx.JSON(http.StatusTooManyRequests, &response.Status{
				Status_Code: -1,
				Status_Msg:  "请求过频",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
