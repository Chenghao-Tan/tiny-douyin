package request

type MessageReq struct {
	Token       string `json:"token" form:"token" binding:"required,jwt"`                     // 用户鉴权token
	To_User_ID  uint   `json:"to_user_id" form:"to_user_id" binding:"required,min=1"`         // 对方用户id
	Action_Type int    `json:"action_type" form:"action_type" binding:"required,min=1,max=1"` // 1-发送消息
	Content     string `json:"content" form:"content" binding:"required,min=1,max=256"`       // 消息内容
}

type MessageListReq struct {
	Token        string `json:"token" form:"token" binding:"required,jwt"`                  // 用户鉴权token
	To_User_ID   uint   `json:"to_user_id" form:"to_user_id" binding:"required,min=1"`      // 对方用户id
	Pre_Msg_Time int64  `json:"pre_msg_time" form:"pre_msg_time" binding:"omitempty,min=0"` // 可选参数，上次最新消息的时间 API文档有误 应有此项且为可选参数
}
