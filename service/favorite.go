package service

import (
	"douyin/repo"
	"douyin/repo/db"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

// 点赞/取消赞
func Favorite(ctx *gin.Context, req *request.FavoriteReq) (resp *response.FavoriteResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("req_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (req_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 存储点赞信息
	if req.Action_Type == 1 {
		// 点赞
		err = repo.CreateUserFavorites(context.TODO(), req_id.(uint), req.Video_ID)
		if err != nil {
			utility.Logger().Errorf("CreateUserFavorites err: %v", err)
			return nil, err
		}
	} else if req.Action_Type == 2 {
		// 取消赞
		err = repo.DeleteUserFavorites(context.TODO(), req_id.(uint), req.Video_ID)
		if err != nil {
			utility.Logger().Errorf("DeleteUserFavorites err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", req.Action_Type)
		return nil, errors.New("操作类型有误")
	}

	return &response.FavoriteResp{}, nil
}

// 获取喜欢列表
func FavoriteList(ctx *gin.Context, req *request.FavoriteListReq) (resp *response.FavoriteListResp, err error) {
	// 读取目标用户信息
	favorites, err := db.ReadUserFavorites(context.TODO(), req.User_ID)
	if err != nil {
		utility.Logger().Errorf("ReadUserFavorites err: %v", err)
		return nil, err
	}

	// 读取目标用户喜欢列表
	resp = &response.FavoriteListResp{Video_List: make([]response.Video, 0, len(favorites))} // 初始化响应
	for _, video := range favorites {
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
