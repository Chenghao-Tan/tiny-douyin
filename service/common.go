package service

import (
	"douyin/repo/db"
	"douyin/repo/oss"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 读取指定用户信息 返回用户信息响应结构体
func readUserInfo(ctx *gin.Context, user_id uint) (userInfo *response.User, err error) {
	// 获取请求用户ID
	req_id, _ := ctx.Get("user_id") // 允许无法获取 获取请求用户ID不成功时req_id为nil

	// 读取目标用户信息 //TODO
	user, err := db.FindUserByID(context.TODO(), user_id)
	if err != nil {
		utility.Logger().Errorf("FindUserByID err: %v", err)
		return nil, err
	}
	followCount := uint(len(user.Follows))     // 统计关注数
	followerCount := uint(len(user.Followers)) // 统计粉丝数
	workCount := uint(len(user.Works))         // 统计作品数
	favoriteCount := uint(len(user.Favorites)) // 统计点赞数

	// 统计获赞数
	var favoritedCount uint = 0
	for _, v := range user.Works {
		favoritedCount += uint(len(v.Favorited))
	}

	// 检查是否被请求用户关注
	isFollow := false
	if req_id != nil {
		isFollow = db.CheckFollow(context.TODO(), req_id.(uint), uint(user.ID))
	}

	// 获取头像及个人页背景图URL
	avatarURL, _ := oss.GetAvatar(context.TODO(), strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utility.Logger().Errorf("GetAvatar err: %v", err) // 允许无法获取 仅记录错误
	}
	backgroundImageURL, _ := oss.GetBackgroundImage(context.TODO(), strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		utility.Logger().Errorf("GetBackgroundImage err: %v", err) // 允许无法获取 仅记录错误
	}

	return &response.User{
		ID:               user.ID,
		Name:             user.Username,
		Follow_Count:     followCount,
		Follower_Count:   followerCount,
		Is_Follow:        isFollow,
		Avatar:           avatarURL,
		Background_Image: backgroundImageURL,
		Signature:        user.Signature,
		Total_Favorited:  strconv.FormatUint(uint64(favoritedCount), 10),
		Work_Count:       workCount,
		Favorite_Count:   favoriteCount,
	}, nil
}

// 读取指定视频信息 返回视频信息响应结构体
func readVideoInfo(ctx *gin.Context, video_id uint) (videoInfo *response.Video, err error) {
	// 获取请求用户ID
	req_id, _ := ctx.Get("user_id") // 允许无法获取 获取请求用户ID不成功时req_id为nil

	// 读取目标视频信息 //TODO
	video, err := db.FindVideoByID(context.TODO(), video_id)
	if err != nil {
		utility.Logger().Errorf("FindVideoByID err: %v", err)
		return nil, err
	}
	favoritedCount := uint(len(video.Favorited)) // 统计获赞数
	commentCount := uint(len(video.Comments))    // 统计评论数

	// 获取视频及封面URL
	videoURL, coverURL, err := oss.GetVideo(context.TODO(), strconv.FormatUint(uint64(video.ID), 10))
	if err != nil {
		utility.Logger().Errorf("GetVideo err: %v", err)
		return nil, err
	}

	// 检查是否被请求用户点赞
	isFavorite := false
	if req_id != nil {
		isFavorite = db.CheckFavorite(context.TODO(), req_id.(uint), video.ID)
	}

	// 读取作者信息
	authorInfo, err := readUserInfo(ctx, video.AuthorID)
	if err != nil {
		utility.Logger().Errorf("readUserInfo err: %v", err)
		return nil, err
	}

	return &response.Video{
		ID:             video.ID,
		Author:         *authorInfo,
		Play_URL:       videoURL,
		Cover_URL:      coverURL,
		Favorite_Count: favoritedCount,
		Comment_Count:  commentCount,
		Is_Favorite:    isFavorite,
		Title:          video.Title,
	}, nil
}
