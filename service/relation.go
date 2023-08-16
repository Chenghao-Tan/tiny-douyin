package service

import (
	"douyin/repository/dao"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"

	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 关注/取消关注
func Follow(ctx *gin.Context, req *request.FollowReq) (resp *response.FollowResp, err error) {
	// 获取请求用户ID
	Me_ID, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: inaccessible")
		return nil, errors.New("无法获取user_id")
	}

	// 获取目标用户ID
	to_user_id, err := strconv.ParseUint(req.To_User_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 关注/取消关注
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	if action_type == 1 {
		err = dao.CreateFollow(context.TODO(), Me_ID.(uint), uint(to_user_id))
		if err != nil {
			utils.ZapLogger.Errorf("CreateFollow err: %v", err)
			return nil, err
		}
	} else if action_type == 2 {
		err = dao.DeleteFollow(context.TODO(), Me_ID.(uint), uint(to_user_id))
		if err != nil {
			utils.ZapLogger.Errorf("DeleteFollow err: %v", err)
			return nil, err
		}
	} else {
		utils.ZapLogger.Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}

	return &response.FollowResp{}, nil
}

// 获取关注列表
func FollowList(ctx *gin.Context, req *request.FollowListReq) (resp *response.FollowListResp, err error) {
	// 获取目标用户信息
	userID, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := dao.FindUserByID(context.TODO(), uint(userID))
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户关注列表
	resp = &response.FollowListResp{}
	for _, follow := range user.Follows {
		followInfo, err := readUserInfo(ctx, follow.ID)
		if err != nil {
			utils.ZapLogger.Errorf("readUserInfo err: %v", err)
			continue // 跳过本用户
		}

		// 加入响应列表
		resp.User_List = append(resp.User_List, *followInfo)
	}

	return resp, nil
}

// 获取粉丝列表
func FollowerList(ctx *gin.Context, req *request.FollowerListReq) (resp *response.FollowerListResp, err error) {
	// 获取目标用户信息
	userID, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := dao.FindUserByID(context.TODO(), uint(userID))
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户粉丝列表
	resp = &response.FollowerListResp{}
	for _, follower := range user.Followers {
		followerInfo, err := readUserInfo(ctx, follower.ID)
		if err != nil {
			utils.ZapLogger.Errorf("readUserInfo err: %v", err)
			continue // 跳过本用户
		}

		// 加入响应列表
		resp.User_List = append(resp.User_List, *followerInfo)
	}

	return resp, nil
}

// 获取好友列表
func FriendList(ctx *gin.Context, req *request.FriendListReq) (resp *response.FriendListResp, err error) {
	// 获取目标用户信息
	userID, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := dao.FindUserByID(context.TODO(), uint(userID))
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户关注列表
	resp = &response.FriendListResp{}
	for _, follow := range user.Follows {
		// 检查该用户是否也关注了目标用户
		if dao.CheckFollow(context.TODO(), follow.ID, user.ID) {
			// 若互粉则为朋友
			followInfo, err := readUserInfo(ctx, follow.ID)
			if err != nil {
				utils.ZapLogger.Errorf("readUserInfo err: %v", err)
				continue // 跳过本用户
			}

			// 获取上一次消息
			friendUser := response.FriendUser{User: *followInfo}
			outMessage, err1 := dao.FindMessagesBy_From_To_ID(context.TODO(), user.ID, follow.ID, time.Now().Unix(), false, 1)
			inMessage, err2 := dao.FindMessagesBy_From_To_ID(context.TODO(), follow.ID, user.ID, time.Now().Unix(), false, 1)
			if (err1 == nil && err2 == nil) && (len(outMessage) > 0 && len(inMessage) > 0) {
				if outMessage[0].CreatedAt.Unix() > inMessage[0].CreatedAt.Unix() {
					friendUser.Message = outMessage[0].Content
					friendUser.Msg_Type = 1 // 目标用户发送的消息
				} else {
					friendUser.Message = inMessage[0].Content
					friendUser.Msg_Type = 0 // 目标用户接收的消息
				}
			} else if ((err1 == nil) && (len(outMessage) > 0)) && ((err2 != nil) || (len(inMessage) == 0)) {
				friendUser.Message = outMessage[0].Content
				friendUser.Msg_Type = 1 // 目标用户发送的消息
			} else if ((err1 != nil) || (len(outMessage) == 0)) && ((err2 == nil) && (len(inMessage) > 0)) {
				friendUser.Message = inMessage[0].Content
				friendUser.Msg_Type = 0 // 目标用户接收的消息
			} else {
				// friendUser.Message = "" // 默认为不发送
				friendUser.Msg_Type = 2 // 无消息往来时根据API文档强制要求将msgType赋值
			}

			// 加入响应列表
			resp.User_List = append(resp.User_List, friendUser)
		}
	}

	return resp, nil
}
