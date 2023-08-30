package request

type FollowReq struct {
	Token       string `json:"token" form:"token" binding:"required,jwt"`                     // 用户鉴权token
	To_User_ID  uint   `json:"to_user_id" form:"to_user_id" binding:"required,min=1"`         // 对方用户id
	Action_Type int    `json:"action_type" form:"action_type" binding:"required,min=1,max=2"` // 1-关注，2-取消关注
}

type FollowListReq struct {
	User_ID uint   `json:"user_id" form:"user_id" binding:"required,min=1"` // 用户id
	Token   string `json:"token" form:"token" binding:"omitempty,jwt"`      // 用户鉴权token API文档有误 应为可选参数
}

type FollowerListReq struct {
	User_ID uint   `json:"user_id" form:"user_id" binding:"required,min=1"` // 用户id
	Token   string `json:"token" form:"token" binding:"omitempty,jwt"`      // 用户鉴权token API文档有误 应为可选参数
}

type FriendListReq struct {
	User_ID uint   `json:"user_id" form:"user_id" binding:"required,min=1"` // 用户id
	Token   string `json:"token" form:"token" binding:"required,jwt"`       // 用户鉴权token
}
