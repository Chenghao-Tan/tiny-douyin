package response

type MessageResp struct {
	Status
}

type MessageListResp struct {
	Status
	Message_List []Message `json:"message_list"` // 消息列表
}
