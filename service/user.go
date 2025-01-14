package service

import (
	"douyin/midware"
	"douyin/repo"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 自定义错误类型
var ErrorUserExists = errors.New("用户已存在")
var ErrorWrongPassword = errors.New("账号或密码错误")

// 默认个人签名
const signature = "Ad Astra Per Aspera"

// 用户注册
func UserRegister(ctx *gin.Context, req *request.UserRegisterReq) (resp *response.UserRegisterResp, err error) {
	// 校验用户名是否可注册
	if !repo.CheckUserRegister(context.TODO(), req.Username) {
		return nil, ErrorUserExists
	}

	// 存储用户信息
	user, err := repo.CreateUser(context.TODO(), req.Username, req.Password, signature)
	if err != nil {
		utility.Logger().Errorf("CreateUser err: %v", err)
		return nil, err
	}

	// 上传默认头像及个人页背景图
	err = repo.UploadAvatarStream(context.TODO(), strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utility.Logger().Errorf("UploadAvatarStream err: %v", err) // 响应为注册成功 仅记录错误
	}
	err = repo.UploadBackgroundImageStream(context.TODO(), strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utility.Logger().Errorf("UploadBackgroundImageStream err: %v", err) // 响应为注册成功 仅记录错误
	}

	// 注册后生成用户鉴权token(自动登录)
	token, err := midware.GenerateToken(user.ID, user.Username)
	if err != nil {
		utility.Logger().Errorf("GenerateToken err: %v", err)
		token = "" // 响应为注册成功 但将无法自动登录
	}

	return &response.UserRegisterResp{User_ID: user.ID, Token: token}, nil
}

// 用户登录
func UserLogin(ctx *gin.Context, req *request.UserLoginReq) (resp *response.UserLoginResp, err error) {
	// 校验用户名密码组合是否有效
	userID, ok := repo.CheckUserLogin(context.TODO(), req.Username, req.Password)
	if !ok {
		return nil, ErrorWrongPassword
	}

	// 校验成功时生成用户鉴权token
	token, err := midware.GenerateToken(userID, req.Username)
	if err != nil {
		utility.Logger().Errorf("GenerateToken err: %v", err)
		return nil, err
	}

	return &response.UserLoginResp{User_ID: userID, Token: token}, nil
}

// 用户信息
func UserInfo(ctx *gin.Context, req *request.UserInfoReq) (resp *response.UserInfoResp, err error) {
	// 读取目标用户信息
	userInfo, err := readUserInfo(ctx, req.User_ID)
	if err != nil {
		utility.Logger().Errorf("readUserInfo err: %v", err)
		return nil, err
	}

	return &response.UserInfoResp{User: *userInfo}, nil
}
