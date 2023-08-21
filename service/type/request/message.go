package request

type MessageReq struct {
	Token       string `json:"token" form:"token" binding:"required"`             // 用户鉴权token
	To_User_ID  string `json:"to_user_id" form:"to_user_id" binding:"required"`   // 对方用户id
	Action_Type string `json:"action_type" form:"action_type" binding:"required"` // 1-发送消息
	Content     string `json:"content" form:"content" binding:"required"`         // 消息内容
}

type MessageListReq struct {
	Token        string `json:"token" form:"token" binding:"required"`           // 用户鉴权token
	To_User_ID   string `json:"to_user_id" form:"to_user_id" binding:"required"` // 对方用户id
	Pre_Msg_Time uint   `json:"pre_msg_time" form:"pre_msg_time"`                // 可选参数，上次最新消息的时间 API文档有误 应有此项且为可选
}
