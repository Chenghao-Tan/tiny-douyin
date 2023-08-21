package api

import (
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"
	"strconv"
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
	if req.Latest_Time == "" || req.Latest_Time == "0" { // 不存在时字符串为空 由数值转换而来的字符串无效值则为0
		req.Latest_Time = strconv.FormatInt(time.Now().Unix(), 10) // 空值表示当前时间
	} else {
		req.Latest_Time = req.Latest_Time[:len(req.Latest_Time)-3] // API文档有误 请求实为毫秒时间戳 故在此转换
	}

	// token字段
	if req.Token != "" {
		// 解析/校验token (自动验证有效期等)
		claims, err := utility.ParseToken(req.Token)
		if err == nil { // 若成功登录
			// 提取user_id和username
			ctx.Set("user_id", uint(claims["user_id"].(float64))) // token中解析数字默认float64
			ctx.Set("username", claims["username"])
		}
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
