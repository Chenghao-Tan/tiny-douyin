package router

import (
	"douyin/api"
	"douyin/conf"
	"douyin/midware"

	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func NewRouter() *gin.Engine {
	capacity := int64(conf.Cfg().System.Capacity)
	recover := int64(conf.Cfg().System.Recover)

	ginRouter := gin.Default()
	rootAPI := ginRouter.Group("/douyin")
	{
		rootAPI.GET("/ping/", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		rootAPI.GET("/feed", midware.MiddlewareRateLimit(capacity, recover), api.GETFeed) // 应用限流中间件

		userAPI := rootAPI.Group("user")
		{
			userAPI.POST("/register/", api.POSTUserRegister)
			userAPI.POST("/login/", api.POSTUserLogin)
			userAPI.GET("/", midware.MiddlewareAuth(), api.GETUserInfo) // 应用jwt鉴权中间件
		}

		publishAPI := rootAPI.Group("publish")
		{
			publishAPI.POST("/action/", midware.MiddlewareAuth(), api.POSTPublish) // 应用jwt鉴权中间件
			publishAPI.GET("/list/", api.GETPublishList)
		}

		favoriteAPI := rootAPI.Group("favorite")
		{
			favoriteAPI.POST("/action/", midware.MiddlewareAuth(), api.POSTFavorite) // 应用jwt鉴权中间件
			favoriteAPI.GET("/list/", api.GETFavoriteList)
		}

		commentAPI := rootAPI.Group("comment")
		{
			commentAPI.POST("/action/", midware.MiddlewareAuth(), api.POSTComment) // 应用jwt鉴权中间件
			commentAPI.GET("/list/", api.GETCommentList)
		}

		relationAPI := rootAPI.Group("relation")
		{
			relationAPI.POST("/action/", midware.MiddlewareAuth(), api.POSTFollow) // 应用jwt鉴权中间件
			relationAPI.GET("/follow/list/", api.GETFollowList)
			relationAPI.GET("/follower/list/", api.GETFollowerList)
			relationAPI.GET("/friend/list/", midware.MiddlewareAuth(), api.GETFriendList) // 应用jwt鉴权中间件
		}

		messageAPI := rootAPI.Group("message")
		{
			messageAPI.POST("/action/", midware.MiddlewareAuth(), api.POSTMessage) // 应用jwt鉴权中间件
			messageAPI.GET("/chat/", midware.MiddlewareAuth(), api.GETMessageList) // 应用jwt鉴权中间件
		}
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
