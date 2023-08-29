package request

import (
	"mime/multipart"
)

type PublishReq struct {
	Token string                `json:"token" form:"token" binding:"required"` // 用户鉴权token
	Data  *multipart.FileHeader `json:"data" form:"data" binding:"required"`   // 视频数据
	Title string                `json:"title" form:"title" binding:"required"` // 视频标题
}

type PublishListReq struct {
	User_ID uint   `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token   string `json:"token,omitempty" form:"token"`              // 用户鉴权token API文档有误 应为可选参数
}
