package service

import (
	"douyin/repo/db"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 关注/取消关注
func Follow(ctx *gin.Context, req *request.FollowReq) (resp *response.FollowResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("user_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (user_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 读取目标用户ID
	to_user_id, err := strconv.ParseUint(req.To_User_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 关注/取消关注
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	if action_type == 1 {
		// 关注
		err = db.CreateFollow(context.TODO(), req_id.(uint), uint(to_user_id))
		if err != nil {
			utility.Logger().Errorf("CreateFollow err: %v", err)
			return nil, err
		}
	} else if action_type == 2 {
		// 取消关注
		err = db.DeleteFollow(context.TODO(), req_id.(uint), uint(to_user_id))
		if err != nil {
			utility.Logger().Errorf("DeleteFollow err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}

	return &response.FollowResp{}, nil
}

// 获取关注列表
func FollowList(ctx *gin.Context, req *request.FollowListReq) (resp *response.FollowListResp, err error) {
	// 读取目标用户信息
	user_id, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := db.FindUserByID(context.TODO(), uint(user_id))
	if err != nil {
		utility.Logger().Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户关注列表 //TODO
	resp = &response.FollowListResp{}
	for _, follow := range user.Follows {
		// 读取被关注用户信息
		followInfo, err := readUserInfo(ctx, follow.ID)
		if err != nil {
			utility.Logger().Errorf("readUserInfo err: %v", err)
			continue // 跳过该用户
		}

		// 将该用户加入列表
		resp.User_List = append(resp.User_List, *followInfo)
	}

	return resp, nil
}

// 获取粉丝列表
func FollowerList(ctx *gin.Context, req *request.FollowerListReq) (resp *response.FollowerListResp, err error) {
	// 读取目标用户信息
	user_id, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := db.FindUserByID(context.TODO(), uint(user_id))
	if err != nil {
		utility.Logger().Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户粉丝列表 //TODO
	resp = &response.FollowerListResp{}
	for _, follower := range user.Followers {
		// 读取粉丝用户信息
		followerInfo, err := readUserInfo(ctx, follower.ID)
		if err != nil {
			utility.Logger().Errorf("readUserInfo err: %v", err)
			continue // 跳过该用户
		}

		// 将该用户加入列表
		resp.User_List = append(resp.User_List, *followerInfo)
	}

	return resp, nil
}

// 获取好友列表
func FriendList(ctx *gin.Context, req *request.FriendListReq) (resp *response.FriendListResp, err error) {
	// 读取目标用户信息
	user_id, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := db.FindUserByID(context.TODO(), uint(user_id))
	if err != nil {
		utility.Logger().Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户关注列表 //TODO
	resp = &response.FriendListResp{}
	for _, follow := range user.Follows {
		// 检查该用户是否也关注了目标用户
		if db.CheckFollow(context.TODO(), follow.ID, user.ID) {
			// 若互粉则为朋友
			// 读取朋友用户信息
			friendInfo, err := readUserInfo(ctx, follow.ID)
			if err != nil {
				utility.Logger().Errorf("readUserInfo err: %v", err)
				continue // 跳过该用户
			}

			// 初始化朋友用户响应结构
			friendUser := response.FriendUser{User: *friendInfo}

			// 获取上一次消息
			outMessage, err1 := db.FindMessagesBy_From_To_ID(context.TODO(), user.ID, follow.ID, time.Now().Unix(), false, 1) // (目标用户)最新发送消息
			inMessage, err2 := db.FindMessagesBy_From_To_ID(context.TODO(), follow.ID, user.ID, time.Now().Unix(), false, 1)  // (目标用户)最新接收消息
			if (err1 == nil && err2 == nil) && (len(outMessage) > 0 && len(inMessage) > 0) {
				// 皆存在
				if outMessage[0].CreatedAt.Unix() > inMessage[0].CreatedAt.Unix() { // 发送消息较新
					friendUser.Message = outMessage[0].Content
					friendUser.Msg_Type = 1 // 使用目标用户发送的消息
				} else { // 接收消息较新
					friendUser.Message = inMessage[0].Content
					friendUser.Msg_Type = 0 // 使用目标用户接收的消息
				}
			} else if ((err1 == nil) && (len(outMessage) > 0)) && ((err2 != nil) || (len(inMessage) == 0)) { // 发送消息存在且接收消息不存在
				friendUser.Message = outMessage[0].Content
				friendUser.Msg_Type = 1 // 使用目标用户发送的消息
			} else if ((err1 != nil) || (len(outMessage) == 0)) && ((err2 == nil) && (len(inMessage) > 0)) { // 接收消息存在且发送消息不存在
				friendUser.Message = inMessage[0].Content
				friendUser.Msg_Type = 0 // 使用目标用户接收的消息
			} else { // 皆不存在
				// friendUser.Message = "" // 默认为不发送
				friendUser.Msg_Type = 2 // 无消息往来时根据API文档强制要求将msgType赋值
			}

			// 将该朋友用户加入列表
			resp.User_List = append(resp.User_List, friendUser)
		}
	}

	return resp, nil
}
