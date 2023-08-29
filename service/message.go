package service

import (
	"douyin/repo/db"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

func Message(ctx *gin.Context, req *request.MessageReq) (resp *response.MessageResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("req_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (req_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 操作消息
	if req.Action_Type == 1 {
		// 发送消息
		_, err := db.CreateMessage(context.TODO(), req_id.(uint), req.To_User_ID, req.Content)
		if err != nil {
			utility.Logger().Errorf("CreateMessage err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", req.Action_Type)
		return nil, errors.New("操作类型有误")
	}

	return &response.MessageResp{}, nil
}

func MessageList(ctx *gin.Context, req *request.MessageListReq) (resp *response.MessageListResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("req_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (req_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 读取消息列表
	messages, err := db.FindMessagesByCreatedAt(context.TODO(), req_id.(uint), req.To_User_ID, req.Pre_Msg_Time, true, -1) // 查找从某刻起新消息 不限制数量
	if err != nil {
		utility.Logger().Errorf("FindMessagesByCreatedAt err: %v", err)
		return nil, err
	}

	resp = &response.MessageListResp{Message_List: make([]response.Message, 0, len(messages))} // 初始化响应
	for _, message := range messages {
		// 初始化消息响应结构
		messageInfo := response.Message{
			ID:           message.ID,
			To_User_ID:   message.ToUserID,
			From_User_ID: message.FromUserID,
			Content:      message.Content,
			Create_Time:  message.CreatedAt.Unix() * 1000, // 消息发送时间 API文档有误 响应实为毫秒时间戳 故在此转换
			// Create_Time:  message.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 将该消息加入列表
		resp.Message_List = append(resp.Message_List, messageInfo)
	}

	return resp, nil
}
