package api

import (
	"douyin/midware"
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"
	"strconv"

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

	// 检查操作类型
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	}
	if action_type == 1 {
		if req.Comment_Text == "" {
			utility.Logger().Errorf("Invalid comment_text err: invalid")
			ctx.JSON(http.StatusBadRequest, &response.Status{
				Status_Code: -1,
				Status_Msg:  "需要有效comment_text",
			})
			return
		}
	} else if action_type == 2 {
		if req.Comment_ID == "" {
			utility.Logger().Errorf("Invalid comment_id err: invalid")
			ctx.JSON(http.StatusBadRequest, &response.Status{
				Status_Code: -1,
				Status_Msg:  "需要有效comment_id",
			})
			return
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", action_type)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作类型有误",
		})
		return
	}

	// 调用评论/删除评论处理
	resp, err := service.Comment(ctx, req)
	if err != nil {
		utility.Logger().Errorf("Comment err: %v", err)
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

	// 处理可选参数
	// token字段
	if req.Token != "" {
		// 解析/校验token (自动验证有效期等)
		claims, err := midware.ParseToken(req.Token)
		if err == nil { // 若成功登录
			// 提取user_id和username
			ctx.Set("req_id", claims.User_ID)
			ctx.Set("username", claims.Username)
		}
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
