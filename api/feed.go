package api

import (
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GETFeed(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FeedReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 处理可选参数
	// latest_time字段
	if req.Latest_Time == 0 { // 不存在时为0
		req.Latest_Time = time.Now().Unix() // 使用当前时间
	} else {
		req.Latest_Time = req.Latest_Time / 1000 // API文档有误 请求实为毫秒时间戳 故在此转换
	}

	// 调用获取视频列表
	resp, err := service.Feed(ctx, req)
	if err != nil {
		utility.Logger().Errorf("Feed err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 获取成功
	status := response.Status{Status_Code: 0, Status_Msg: "获取成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}
