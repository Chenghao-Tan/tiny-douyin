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

func POSTMessage(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.MessageReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	}

	// 调用消息发送处理
	resp, err := service.Message(ctx, req)
	if err != nil {
		utility.Logger().Errorf("Message err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	}

	// 操作成功
	status := response.Status{Status_Code: 0, Status_Msg: "操作成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}

func GETMessageList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.MessageListReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 处理特殊参数
	// pre_msg_time字段 不存在时req.Pre_Msg_Time为0 此时同样适用于以下处理
	req.Pre_Msg_Time = req.Pre_Msg_Time / 1000 // API文档有误 请求实为毫秒时间戳 故在此转换
	if req.Pre_Msg_Time > time.Now().Unix()+1 {
		// 请求时间晚于当前时间+1秒, 必定无更新消息, 直接响应为获取成功
		status := response.Status{Status_Code: 0, Status_Msg: "获取成功"}
		ctx.JSON(http.StatusOK, status)
		return
	}

	// 调用获取消息记录
	resp, err := service.MessageList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("MessageList err: %v", err)
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
