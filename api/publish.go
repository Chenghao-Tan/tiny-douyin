package api

import (
	"douyin/service"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func POSTPublish(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.PublishReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utils.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "发布失败: " + err.Error(),
		})
		return
	}

	// 调用投稿处理
	resp, err := service.Publish(ctx, req)
	if err != nil {
		utils.Logger().Errorf("Publish err: %v", err)
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
		utils.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	/*// 从请求中读取目标用户ID并与token比对
	user_id, ok := ctx.Get("user_id")
	if !ok || req.User_ID != strconv.FormatUint(uint64(user_id.(uint)), 10) {
		utils.Logger().Errorf("GETPublishList err: 查询目标与请求用户不同")
		ctx.JSON(http.StatusUnauthorized, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}*/

	// 调用获取发布列表
	resp, err := service.PublishList(ctx, req)
	if err != nil {
		utils.Logger().Errorf("PublishList err: %v", err)
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
