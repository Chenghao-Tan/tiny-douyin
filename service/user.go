package service

import (
	"douyin/repository/dao"
	"douyin/repository/model"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"
	"douyin/utils/oss"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 自定义错误类型
var ErrorUserExists = errors.New("用户已存在")
var ErrorUserNotExists = errors.New("用户不存在")
var ErrorWrongPassword = errors.New("账号或密码错误")

// 用户注册
func UserRegister(ctx *gin.Context, req *request.UserRegisterReq) (resp *response.UserRegisterResp, err error) {
	// 查找用户是否已存在 //TODO
	_, err = dao.FindUserByUsername(context.TODO(), req.Username)
	if err == nil {
		err = ErrorUserExists
		return nil, err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) { // 若出现查找功能性错误
		utils.Logger().Errorf("FindUserByUserName err: %v", err)
		return nil, err
	}

	// 只有err == ErrRecordNotFound时才可以注册
	// 准备要存储的内容
	user := &model.User{
		Username:  req.Username,
		Signature: "Ad Astra Per Aspera", // 添加默认签名
	}
	err = user.SetPassword(req.Password) // 向要存储的内容中添加单向加密后的密码
	if err != nil {
		utils.Logger().Errorf("SetPassword err: %v", err)
		return nil, err
	}

	// 存储用户信息 //TODO
	user, err = dao.CreateUser(context.TODO(), user)
	if err != nil {
		utils.Logger().Errorf("CreateUser err: %v", err)
		return nil, err
	}

	// 上传默认头像及个人页背景图
	err = oss.UploadAvatarStream(context.TODO(), strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utils.Logger().Errorf("UploadAvatarStream err: %v", err) // 响应为注册成功 仅记录错误
	}
	err = oss.UploadBackgroundImageStream(context.TODO(), strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utils.Logger().Errorf("UploadBackgroundImageStream err: %v", err) // 响应为注册成功 仅记录错误
	}

	// 注册后生成用户鉴权token(自动登录)
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Logger().Errorf("GenerateToken err: %v", err)
		token = "" // 响应为注册成功 但将无法自动登录
	}

	return &response.UserRegisterResp{User_ID: user.ID, Token: token}, nil
}

// 用户登录
func UserLogin(ctx *gin.Context, req *request.UserLoginReq) (resp *response.UserLoginResp, err error) {
	// 查找用户是否已存在
	user, err := dao.FindUserByUsername(context.TODO(), req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ErrorUserNotExists
			return nil, err
		} else { // 若出现查找功能性错误
			utils.Logger().Errorf("FindUserByUserName err: %v", err)
			return nil, err
		}
	}

	// 校验密码
	if !user.CheckPassword(req.Password) {
		err = ErrorWrongPassword
		return nil, err
	}

	// 校验成功时生成用户鉴权token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Logger().Errorf("GenerateToken err: %v", err)
		return nil, err
	}

	return &response.UserLoginResp{User_ID: user.ID, Token: token}, nil
}

// 用户信息
func UserInfo(ctx *gin.Context, req *request.UserInfoReq) (resp *response.UserInfoResp, err error) {
	// 读取目标用户信息
	user_id, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utils.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	userInfo, err := readUserInfo(ctx, uint(user_id))
	if err != nil {
		utils.Logger().Errorf("readUserInfo err: %v", err)
		return nil, err
	}

	return &response.UserInfoResp{User: *userInfo}, nil
}
