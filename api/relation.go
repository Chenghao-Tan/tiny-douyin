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

func POSTFollow(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FollowReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	}

	// 粗略过滤目标用户ID
	maxID, err := repo.MaxUserID(context.TODO())
	if err == nil {
		if req.To_User_ID > maxID+1 { // 因maxID不一定完全准确, 预留一定余量
			utility.Logger().Warnf("POSTFollow warn: ID越界%v", req.To_User_ID)
			ctx.JSON(http.StatusForbidden, &response.Status{ // 防止泄露maxID
				Status_Code: -1,
				Status_Msg:  "无权操作",
			})
		}
	}

	// 调用关注/取消关注处理
	resp, err := service.Follow(ctx, req)
	if err != nil {
		utility.Logger().Errorf("Follow err: %v", err)
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

func GETFollowList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FollowListReq{}
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
			utility.Logger().Warnf("GETFollowList warn: ID越界%v", req.User_ID)
			ctx.JSON(http.StatusForbidden, &response.Status{ // 防止泄露maxID
				Status_Code: -1,
				Status_Msg:  "无权操作",
			})
		}
	}

	// 调用获取关注列表
	resp, err := service.FollowList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("FollowList err: %v", err)
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

func GETFollowerList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FollowerListReq{}
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
			utility.Logger().Warnf("GETFollowerList warn: ID越界%v", req.User_ID)
			ctx.JSON(http.StatusForbidden, &response.Status{ // 防止泄露maxID
				Status_Code: -1,
				Status_Msg:  "无权操作",
			})
		}
	}

	// 调用获取粉丝列表
	resp, err := service.FollowerList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("FollowerList err: %v", err)
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

func GETFriendList(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.FriendListReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "获取失败: " + err.Error(),
		})
		return
	}

	// 从请求中读取目标用户ID并与token比对
	req_id, ok := ctx.Get("req_id")
	if !ok || req.User_ID != req_id.(uint) {
		utility.Logger().Warnf("GETFriendList warn: 查询目标与请求用户不同")
		ctx.JSON(http.StatusForbidden, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}

	// 粗略过滤目标用户ID
	maxID, err := repo.MaxUserID(context.TODO())
	if err == nil {
		if req.User_ID > maxID+1 { // 因maxID不一定完全准确, 预留一定余量
			utility.Logger().Warnf("GETFriendList warn: ID越界%v", req.User_ID)
			ctx.JSON(http.StatusForbidden, &response.Status{ // 防止泄露maxID
				Status_Code: -1,
				Status_Msg:  "无权操作",
			})
		}
	}

	// 调用获取好友列表
	resp, err := service.FriendList(ctx, req)
	if err != nil {
		utility.Logger().Errorf("FriendList err: %v", err)
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
