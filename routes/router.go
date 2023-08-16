package routes

import (
	"douyin/api"
	"douyin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	rootApi := ginRouter.Group("/douyin")
	{
		rootApi.GET("/ping/", func(context *gin.Context) {
			context.JSON(http.StatusOK, "success")
		})

		rootApi.GET("/feed/", utils.MiddlewareRateLimit(10, 1), api.GETFeed) // 应用限流中间件 最大10次/秒 每秒恢复1次

		userApi := rootApi.Group("user")
		{
			userApi.POST("/register/", api.POSTUserRegister)
			userApi.POST("/login/", api.POSTUserLogin)
			userApi.GET("/", utils.MiddlewareAuth(), api.GETUserInfo) // 应用jwt鉴权中间件
		}

		publishApi := rootApi.Group("publish")
		{
			publishApi.POST("/action/", utils.MiddlewareAuth(), api.POSTPublish) // 应用jwt鉴权中间件
			publishApi.GET("/list/", utils.MiddlewareAuth(), api.GETPublishList) // 应用jwt鉴权中间件
		}
	}
	return ginRouter
}
