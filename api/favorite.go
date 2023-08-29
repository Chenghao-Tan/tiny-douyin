package api

import (
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"

	"github.com/gin-gonic/gin"
)

func POSTFavorite(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FavoriteReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	}

	// 检查操作类型
	if !(req.Action_Type == 1 || req.Action_Type == 2) {
		utility.Logger().Errorf("Invalid action_type err: %v", req.Action_Type)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作类型有误",
		})
		return
	}

	// 调用赞/取消赞处理
	resp, err := service.Favorite(ctx, req)
	if err != nil {
		utility.Logger().Errorf("Favorite err: %v", err)
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

func GETFavoriteList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FavoriteListReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 调用获取喜欢列表
	resp, err := service.FavoriteList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("FavoriteList err: %v", err)
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
