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
func readUserInfo(ctx *gin.Context, userID uint) (userInfo *response.User, err error) {
	// 获取请求用户ID
	req_id, _ := ctx.Get("req_id") // 允许无法获取 获取请求用户ID不成功时req_id为nil

	// 读取目标用户信息
	user, err := db.FindUserByID(context.TODO(), userID)
	if err != nil {
		utility.Logger().Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	followCount := uint(db.CountUserFollows(context.TODO(), userID))      // 统计关注数
	followerCount := uint(db.CountUserFollowers(context.TODO(), userID))  // 统计粉丝数
	workCount := uint(db.CountUserWorks(context.TODO(), userID))          // 统计作品数
	favoriteCount := uint(db.CountUserFavorites(context.TODO(), userID))  // 统计点赞数
	favoritedCount := uint(db.CountUserFavorited(context.TODO(), userID)) // 统计获赞数

	// 检查是否被请求用户关注
	isFollow := false
	if req_id != nil {
		isFollow = db.CheckFollow(context.TODO(), req_id.(uint), uint(userID))
	}

	// 获取头像及个人页背景图URL
	avatarURL, _ := oss.GetAvatar(context.TODO(), strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		utility.Logger().Errorf("GetAvatar err: %v", err) // 允许无法获取 仅记录错误
	}
	backgroundImageURL, _ := oss.GetBackgroundImage(context.TODO(), strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		utility.Logger().Errorf("GetBackgroundImage err: %v", err) // 允许无法获取 仅记录错误
	}

	return &response.User{
		ID:               userID,
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
func readVideoInfo(ctx *gin.Context, videoID uint) (videoInfo *response.Video, err error) {
	// 获取请求用户ID
	req_id, _ := ctx.Get("req_id") // 允许无法获取 获取请求用户ID不成功时req_id为nil

	// 读取目标视频信息
	video, err := db.FindVideoByID(context.TODO(), videoID)
	if err != nil {
		utility.Logger().Errorf("FindVideoByID err: %v", err)
		return nil, err
	}

	favoritedCount := uint(db.CountVideoFavorited(context.TODO(), videoID)) // 统计获赞数
	commentCount := uint(db.CountVideoComments(context.TODO(), videoID))    // 统计评论数

	// 获取视频及封面URL
	videoURL, coverURL, err := oss.GetVideo(context.TODO(), strconv.FormatUint(uint64(videoID), 10))
	if err != nil {
		utility.Logger().Errorf("GetVideo err: %v", err)
		return nil, err
	}

	// 检查是否被请求用户点赞
	isFavorite := false
	if req_id != nil {
		isFavorite = db.CheckFavorite(context.TODO(), req_id.(uint), videoID)
	}

	// 读取作者信息
	authorInfo, err := readUserInfo(ctx, video.AuthorID)
	if err != nil {
		utility.Logger().Errorf("readUserInfo err: %v", err)
		return nil, err
	}

	return &response.Video{
		ID:             videoID,
		Author:         *authorInfo,
		Play_URL:       videoURL,
		Cover_URL:      coverURL,
		Favorite_Count: favoritedCount,
		Comment_Count:  commentCount,
		Is_Favorite:    isFavorite,
		Title:          video.Title,
	}, nil
}
