package router

import (
	"douyin/api"
	"douyin/conf"
	"douyin/midware"

	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func NewRouter() *gin.Engine {
	rateLimit := conf.Cfg().System.RateLimit

	ginRouter := gin.Default()
	rootAPI := ginRouter.Group("/douyin")
	{
		rootAPI.GET("/ping/", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		rootAPI.GET("/feed", midware.MiddlewareRateLimitWithRedis(rateLimit, time.Second), midware.MiddlewareAuth(false), api.GETFeed) // 应用限流中间件, jwt鉴权中间件

		userAPI := rootAPI.Group("user")
		{
			userAPI.POST("/register/", midware.MiddlewareRateLimitWithRedis(rateLimit, time.Second), api.POSTUserRegister) // 应用限流中间件
			userAPI.POST("/login/", midware.MiddlewareRateLimitWithRedis(rateLimit, time.Second), api.POSTUserLogin)       // 应用限流中间件
			userAPI.GET("/", midware.MiddlewareAuth(true), api.GETUserInfo)                                                // 应用jwt鉴权中间件(强制)
		}

		publishAPI := rootAPI.Group("publish")
		{
			publishAPI.POST("/action/", midware.MiddlewareAuth(true), api.POSTPublish)  // 应用jwt鉴权中间件(强制)
			publishAPI.GET("/list/", midware.MiddlewareAuth(false), api.GETPublishList) // 应用jwt鉴权中间件
		}

		favoriteAPI := rootAPI.Group("favorite")
		{
			favoriteAPI.POST("/action/", midware.MiddlewareAuth(true), api.POSTFavorite)  // 应用jwt鉴权中间件(强制)
			favoriteAPI.GET("/list/", midware.MiddlewareAuth(false), api.GETFavoriteList) // 应用jwt鉴权中间件
		}

		commentAPI := rootAPI.Group("comment")
		{
			commentAPI.POST("/action/", midware.MiddlewareAuth(true), api.POSTComment)  // 应用jwt鉴权中间件(强制)
			commentAPI.GET("/list/", midware.MiddlewareAuth(false), api.GETCommentList) // 应用jwt鉴权中间件
		}

		relationAPI := rootAPI.Group("relation")
		{
			relationAPI.POST("/action/", midware.MiddlewareAuth(true), api.POSTFollow)             // 应用jwt鉴权中间件(强制)
			relationAPI.GET("/follow/list/", midware.MiddlewareAuth(false), api.GETFollowList)     // 应用jwt鉴权中间件
			relationAPI.GET("/follower/list/", midware.MiddlewareAuth(false), api.GETFollowerList) // 应用jwt鉴权中间件
			relationAPI.GET("/friend/list/", midware.MiddlewareAuth(true), api.GETFriendList)      // 应用jwt鉴权中间件(强制)
		}

		messageAPI := rootAPI.Group("message")
		{
			messageAPI.POST("/action/", midware.MiddlewareAuth(true), api.POSTMessage) // 应用jwt鉴权中间件(强制)
			messageAPI.GET("/chat/", midware.MiddlewareAuth(true), api.GETMessageList) // 应用jwt鉴权中间件(强制)
		}
	}

	var err error
	if strings.ToLower(conf.Cfg().System.TrustedProxy) == "none" {
		err = ginRouter.SetTrustedProxies(nil)
	} else {
		err = ginRouter.SetTrustedProxies([]string{conf.Cfg().System.TrustedProxy})
	}
	if err != nil {
		panic(err)
	}

	return ginRouter
}

func RunWithContext(ctx context.Context, r http.Handler, addr string) (err error) {
	var g errgroup.Group

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	g.Go(func() error {
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done() // 阻塞等待终止信号
		return srv.Shutdown(context.Background())
	})

	return g.Wait()
}
