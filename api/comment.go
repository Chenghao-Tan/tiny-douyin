package api

import (
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"

	"github.com/gin-gonic/gin"
)

func POSTComment(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.CommentReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	}

	// 调用评论/删除评论处理
	resp, err := service.Comment(ctx, req)
	if err != nil {
		var httpCode int
		if err == service.ErrorCommentInaccessible {
			utility.Logger().Warnf("Comment warn: %v", err)
			httpCode = http.StatusForbidden
		} else {
			utility.Logger().Errorf("Comment err: %v", err)
			httpCode = http.StatusInternalServerError
		}
		ctx.JSON(httpCode, &response.Status{
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

func GETCommentList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.CommentListReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 调用获取评论列表
	resp, err := service.CommentList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("CommentList err: %v", err)
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
