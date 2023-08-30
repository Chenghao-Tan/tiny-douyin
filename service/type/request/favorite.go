package request

type FavoriteReq struct {
	Token       string `json:"token" form:"token" binding:"required,jwt"`                     // 用户鉴权token
	Video_ID    uint   `json:"video_id" form:"video_id" binding:"required,min=1"`             // 视频id
	Action_Type int    `json:"action_type" form:"action_type" binding:"required,min=1,max=2"` // 1-点赞，2-取消点赞
}

type FavoriteListReq struct {
	User_ID uint   `json:"user_id" form:"user_id" binding:"required,min=1"` // 用户id
	Token   string `json:"token" form:"token" binding:"omitempty,jwt"`      // 用户鉴权token API文档有误 应为可选参数
}
