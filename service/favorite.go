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

// 点赞/取消赞
func Favorite(ctx *gin.Context, req *request.FavoriteReq) (resp *response.FavoriteResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("user_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (user_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 读取目标视频ID
	video_id, err := strconv.ParseUint(req.Video_ID, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}

	// 存储点赞信息
	action_type, err := strconv.ParseUint(req.Action_Type, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseUint err: %v", err)
		return nil, err
	}
	if action_type == 1 {
		// 点赞
		err = db.CreateFavorite(context.TODO(), req_id.(uint), uint(video_id))
		if err != nil {
			utility.Logger().Errorf("CreateFavorite err: %v", err)
			return nil, err
		}
	} else if action_type == 2 {
		// 取消赞
		err = db.DeleteFavorite(context.TODO(), req_id.(uint), uint(video_id))
		if err != nil {
			utility.Logger().Errorf("DeleteFavorite err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", action_type)
		return nil, errors.New("操作类型有误")
	}

	return &response.FavoriteResp{}, nil
}

// 获取喜欢列表
func FavoriteList(ctx *gin.Context, req *request.FavoriteListReq) (resp *response.FavoriteListResp, err error) {
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

	// 读取目标用户喜欢列表 //TODO
	resp = &response.FavoriteListResp{} // 初始化响应
	for _, video := range user.Favorites {
		// 读取视频信息
		videoInfo, err := readVideoInfo(ctx, video.ID)
		if err != nil {
			utility.Logger().Errorf("readVideoInfo err: %v", err)
			continue // 跳过本条视频
		}

		// 将该视频加入列表
		resp.Video_List = append(resp.Video_List, *videoInfo)
	}

	return resp, nil
}
