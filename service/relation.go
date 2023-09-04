package service

import (
	"douyin/repo"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

// 关注/取消关注
func Follow(ctx *gin.Context, req *request.FollowReq) (resp *response.FollowResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("req_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (req_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 关注/取消关注
	if req.Action_Type == 1 {
		// 关注
		err = repo.CreateUserFollows(context.TODO(), req_id.(uint), req.To_User_ID)
		if err != nil {
			utility.Logger().Errorf("CreateUserFollows err: %v", err)
			return nil, err
		}
	} else if req.Action_Type == 2 {
		// 取消关注
		err = repo.DeleteUserFollows(context.TODO(), req_id.(uint), req.To_User_ID)
		if err != nil {
			utility.Logger().Errorf("DeleteUserFollows err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", req.Action_Type)
		return nil, errors.New("操作类型有误")
	}

	return &response.FollowResp{}, nil
}

// 获取关注列表
func FollowList(ctx *gin.Context, req *request.FollowListReq) (resp *response.FollowListResp, err error) {
	// 读取目标用户信息
	follows, err := repo.ReadUserFollows(context.TODO(), req.User_ID)
	if err != nil {
		utility.Logger().Errorf("ReadUserFollows err: %v", err)
		return nil, err
	}

	// 读取目标用户关注列表
	resp = &response.FollowListResp{User_List: make([]response.User, 0, len(follows))} // 初始化响应
	for _, follow := range follows {
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
	followers, err := repo.ReadUserFollowers(context.TODO(), req.User_ID)
	if err != nil {
		utility.Logger().Errorf("ReadUserFollowers err: %v", err)
		return nil, err
	}

	// 读取目标用户粉丝列表
	resp = &response.FollowerListResp{User_List: make([]response.User, 0, len(followers))} // 初始化响应
	for _, follower := range followers {
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
	follows, err := repo.ReadUserFollows(context.TODO(), req.User_ID)
	if err != nil {
		utility.Logger().Errorf("ReadUserFollows err: %v", err)
		return nil, err
	}

	// 读取目标用户关注列表(用于读取朋友)
	resp = &response.FriendListResp{} // 初始化响应 由于朋友(互粉)数未知且一般较小, 不预先分配空间
	for _, friend := range follows {
		// 检查该用户是否也关注了目标用户
		if repo.CheckUserFollows(context.TODO(), friend.ID, req.User_ID) {
			// 若互粉则为朋友
			// 读取朋友用户信息
			friendInfo, err := readUserInfo(ctx, friend.ID)
			if err != nil {
				utility.Logger().Errorf("readUserInfo err: %v", err)
				continue // 跳过该用户
			}

			// 初始化朋友用户增补响应结构
			friendUser := response.FriendUser{User: *friendInfo}

			// 查找最近一条消息
			message, err := repo.FindMessagesByCreatedAt(context.TODO(), req.User_ID, friend.ID, time.Now().Unix(), false, 1)
			if err != nil {
				utility.Logger().Errorf("FindMessagesByCreatedAt err: %v", err)
				// 响应为获取成功 但最近消息将为空
				// friendUser.Message = "" // 和该好友的最新聊天消息 根据API文档默认为不发送
				friendUser.Msg_Type = 2 // 无消息往来时根据API文档强制要求将msgType赋值
			} else if len(message) == 0 {
				// 最近消息为空
				// friendUser.Message = "" // 和该好友的最新聊天消息 根据API文档默认为不发送
				friendUser.Msg_Type = 2 // 无消息往来时根据API文档强制要求将msgType赋值
			} else if message[0].FromUserID == req.User_ID { // 为目标用户发送的消息
				friendUser.Message = message[0].Content
				friendUser.Msg_Type = 1
			} else if message[0].ToUserID == req.User_ID { // 为目标用户接收的消息
				friendUser.Message = message[0].Content
				friendUser.Msg_Type = 0
			} else {
				utility.Logger().Errorf("FindMessagesByCreatedAt err: 查找结果错误")
				// 响应为获取成功 但最近消息将为空
				// friendUser.Message = "" // 和该好友的最新聊天消息 根据API文档默认为不发送
				friendUser.Msg_Type = 2 // 无消息往来时根据API文档强制要求将msgType赋值
			}

			// 将该朋友用户加入列表
			resp.User_List = append(resp.User_List, friendUser)
		}
	}

	return resp, nil
}
