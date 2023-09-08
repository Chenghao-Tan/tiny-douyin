package api

import (
	"douyin/repo"
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func POSTPublish(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.PublishReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "发布失败: " + err.Error(),
		})
		return
	}

	// 调用投稿处理
	resp, err := service.Publish(ctx, req)
	if err != nil {
		utility.Logger().Errorf("Publish err: %v", err)
		ctx.JSON(http.StatusInternalServerError, &response.Status{
			Status_Code: -1,
			Status_Msg:  "发布失败: " + err.Error(),
		})
		return
	}

	// 发布成功
	status := response.Status{Status_Code: 0, Status_Msg: "发布成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}

func GETPublishList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.PublishListReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 粗略过滤目标用户ID
	maxID, err := repo.MaxUserID(context.TODO())
	if err == nil {
		if req.User_ID > maxID+1 { // 因maxID不一定完全准确, 预留一定余量
			utility.Logger().Warnf("GETPublishList warn: ID越界%v", req.User_ID)
			ctx.JSON(http.StatusForbidden, &response.Status{ // 防止泄露maxID
				Status_Code: -1,
				Status_Msg:  "无权操作",
			})
		}
	}

	// 调用获取发布列表
	resp, err := service.PublishList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("PublishList err: %v", err)
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
