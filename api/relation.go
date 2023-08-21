package api

import (
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"
	"strconv"

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

	// 检查操作类型
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作失败: " + err.Error(),
		})
		return
	} else if !(action_type == 1 || action_type == 2) {
		utility.Logger().Errorf("Invalid action_type err: %v", action_type)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "操作类型有误",
		})
		return
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

	// 处理可选参数
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

	/*// 从请求中读取目标用户ID并与token比对
	user_id, ok := ctx.Get("user_id")
	if !ok || req.User_ID != strconv.FormatUint(uint64(user_id.(uint)), 10) {
		utility.Logger().Errorf("GETFollowList err: 查询目标与请求用户不同")
		ctx.JSON(http.StatusUnauthorized, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}*/

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

	// 处理可选参数
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

	/*// 从请求中读取目标用户ID并与token比对
	user_id, ok := ctx.Get("user_id")
	if !ok || req.User_ID != strconv.FormatUint(uint64(user_id.(uint)), 10) {
		utility.Logger().Errorf("GETFollowerList err: 查询目标与请求用户不同")
		ctx.JSON(http.StatusUnauthorized, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}*/

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

	/*// 从请求中读取目标用户ID并与token比对
	user_id, ok := ctx.Get("user_id")
	if !ok || req.User_ID != strconv.FormatUint(uint64(user_id.(uint)), 10) {
		utility.Logger().Errorf("GETFriendList err: 查询目标与请求用户不同")
		ctx.JSON(http.StatusUnauthorized, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}*/

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
