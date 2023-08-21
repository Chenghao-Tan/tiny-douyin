package request

type FavoriteReq struct {
	Token       string `json:"token" form:"token" binding:"required"`             // 用户鉴权token
	Video_ID    string `json:"video_id" form:"video_id" binding:"required"`       // 视频id
	Action_Type string `json:"action_type" form:"action_type" binding:"required"` // 1-点赞，2-取消点赞
}

type FavoriteListReq struct {
	User_ID string `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token   string `json:"token" form:"token" binding:"required"`     // 用户鉴权token
}
