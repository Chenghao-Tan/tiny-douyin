package service

import (
	"douyin/repo/db"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Message(ctx *gin.Context, req *request.MessageReq) (resp *response.MessageResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("user_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (user_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 获取目标用户ID
	to_user_id, err := strconv.ParseUint(req.To_User_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 操作消息
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	if action_type == 1 {
		// 发送消息 //TODO
		_, err := db.CreateMessage(context.TODO(), req_id.(uint), uint(to_user_id), req.Content)
		if err != nil {
			utility.Logger().Errorf("CreateMessage err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}

	return &response.MessageResp{}, nil
}

func MessageList(ctx *gin.Context, req *request.MessageListReq) (resp *response.MessageListResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("user_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (user_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 获取目标用户ID
	to_user_id, err := strconv.ParseUint(req.To_User_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 读取消息列表
	outMessages, err1 := db.FindMessagesBy_From_To_ID(context.TODO(), req_id.(uint), uint(to_user_id), int64(req.Pre_Msg_Time), true, -1) // 查找从某刻起目标用户发送新消息 不限制数量
	if err1 != nil {
		utility.Logger().Errorf("FindMessagesBy_From_To_ID err: %v", err)
		return nil, err1
	}
	inMessages, err2 := db.FindMessagesBy_From_To_ID(context.TODO(), uint(to_user_id), req_id.(uint), int64(req.Pre_Msg_Time), true, -1) // 查找从某刻起目标用户接收新消息 不限制数量
	if err2 != nil {
		utility.Logger().Errorf("FindMessagesBy_From_To_ID err: %v", err)
		return nil, err2
	}
	messages := append(outMessages, inMessages...) // 合并二者

	resp = &response.MessageListResp{} // 初始化响应
	for _, message := range messages {
		// 初始化消息响应结构
		messageInfo := response.Message{
			ID:           message.ID,
			To_User_ID:   message.ToUserID,
			From_User_ID: message.FromUserID,
			Content:      message.Content,
			Create_Time:  uint(message.CreatedAt.Unix() * 1000), // 消息发送时间 API文档有误 响应实为毫秒时间戳 故在此转换
			// Create_Time:  message.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 将该消息加入列表
		resp.Message_List = append(resp.Message_List, messageInfo)
	}

	return resp, nil
}
