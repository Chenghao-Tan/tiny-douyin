package api

import (
	"douyin/service"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"net/http"

	"github.com/gin-gonic/gin"
)

func POSTUserRegister(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserRegisterReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "注册失败: " + err.Error(),
		})
		return
	}

	// 调用用户注册处理
	resp, err := service.UserRegister(ctx, req)
	if err != nil {
		var httpCode int
		if err == service.ErrorUserExists {
			utility.Logger().Warnf("UserRegister warn: %v", err)
			httpCode = http.StatusConflict
		} else {
			utility.Logger().Errorf("UserRegister err: %v", err)
			httpCode = http.StatusInternalServerError
		}
		ctx.JSON(httpCode, &response.Status{
			Status_Code: -1,
			Status_Msg:  "注册失败: " + err.Error(),
		})
		return
	}

	// 注册成功
	status := response.Status{Status_Code: 0, Status_Msg: "注册成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}

func POSTUserLogin(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserLoginReq{}
	err := ctx.ShouldBind(req)
	if err != nil {
		utility.Logger().Errorf("ShouldBind err: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Status{
			Status_Code: -1,
			Status_Msg:  "登录失败: " + err.Error(),
		})
		return
	}

	// 调用用户登录处理
	resp, err := service.UserLogin(ctx, req)
	if err != nil {
		var httpCode int
		if err == service.ErrorWrongPassword {
			utility.Logger().Warnf("UserLogin warn: %v", err)
			httpCode = http.StatusUnauthorized
		} else {
			utility.Logger().Errorf("UserLogin err: %v", err)
			httpCode = http.StatusInternalServerError
		}
		ctx.JSON(httpCode, &response.Status{
			Status_Code: -1,
			Status_Msg:  "登录失败: " + err.Error(),
		})
		return
	}

	// 登录成功
	status := response.Status{Status_Code: 0, Status_Msg: "登录成功"}
	resp.Status = status
	ctx.JSON(http.StatusOK, resp)
}

func GETUserInfo(ctx *gin.Context) {
	// 绑定JSON到结构体
	req := &request.UserInfoReq{}
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
		utility.Logger().Warnf("GETUserInfo warn: 查询目标与请求用户不同")
		ctx.JSON(http.StatusForbidden, &response.Status{
			Status_Code: -1,
			Status_Msg:  "无权获取",
		})
		return
	}

	// 调用获取用户信息
	resp, err := service.UserInfo(ctx, req)
	if err != nil {
		utility.Logger().Errorf("UserInfo err: %v", err)
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
