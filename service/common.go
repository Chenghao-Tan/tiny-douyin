package service

import (
	"douyin/repo"
	"douyin/repo/db"
	"douyin/repo/oss"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 读取指定用户信息 返回用户信息响应结构体
func readUserInfo(ctx *gin.Context, userID uint) (userInfo *response.User, err error) {
	// 获取请求用户ID
	req_id, _ := ctx.Get("req_id") // 允许无法获取 获取请求用户ID不成功时req_id为nil

	// 读取目标用户基本信息
	user, err := db.ReadUserBasics(context.TODO(), userID)
	if err != nil {
		utility.Logger().Errorf("ReadUserBasics err: %v", err)
		return nil, err
	}

	followCount := uint(db.CountUserFollows(context.TODO(), userID))        // 统计关注数
	followerCount := uint(db.CountUserFollowers(context.TODO(), userID))    // 统计粉丝数
	workCount := uint(db.CountUserWorks(context.TODO(), userID))            // 统计作品数
	favoriteCount := uint(repo.CountUserFavorites(context.TODO(), userID))  // 统计点赞数
	favoritedCount := uint(repo.CountUserFavorited(context.TODO(), userID)) // 统计获赞数

	// 检查是否被请求用户关注
	isFollow := false
	if req_id != nil {
		isFollow = db.CheckUserFollows(context.TODO(), req_id.(uint), uint(userID))
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
		Total_Favorited:  favoritedCount,
		Work_Count:       workCount,
		Favorite_Count:   favoriteCount,
	}, nil
}

// 读取指定视频信息 返回视频信息响应结构体
func readVideoInfo(ctx *gin.Context, videoID uint) (videoInfo *response.Video, err error) {
	// 获取请求用户ID
	req_id, _ := ctx.Get("req_id") // 允许无法获取 获取请求用户ID不成功时req_id为nil

	// 读取目标视频基本信息
	video, err := db.ReadVideoBasics(context.TODO(), videoID)
	if err != nil {
		utility.Logger().Errorf("ReadVideoBasics err: %v", err)
		return nil, err
	}

	favoritedCount := uint(repo.CountVideoFavorited(context.TODO(), videoID)) // 统计获赞数
	commentCount := uint(db.CountVideoComments(context.TODO(), videoID))      // 统计评论数

	// 获取视频及封面URL
	videoURL, coverURL, err := oss.GetVideo(context.TODO(), strconv.FormatUint(uint64(videoID), 10))
	if err != nil {
		utility.Logger().Errorf("GetVideo err: %v", err)
		return nil, err
	}

	// 检查是否被请求用户点赞
	isFavorite := false
	if req_id != nil {
		isFavorite = repo.CheckUserFavorites(context.TODO(), req_id.(uint), videoID)
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

// 读取指定评论信息 返回评论信息响应结构体
func readCommentInfo(ctx *gin.Context, commentID uint) (commentInfo *response.Comment, err error) {
	// 读取目标评论基本信息
	comment, err := db.ReadCommentBasics(context.TODO(), commentID)
	if err != nil {
		utility.Logger().Errorf("ReadCommentBasics err: %v", err)
		return nil, err
	}

	// 读取作者信息
	authorInfo, err := readUserInfo(ctx, comment.AuthorID)
	if err != nil {
		utility.Logger().Errorf("readUserInfo err: %v", err)
		return nil, err
	}

	return &response.Comment{
		ID:          commentID,
		User:        *authorInfo,
		Content:     comment.Content,
		Create_Date: fmt.Sprintf("%02d-%02d", comment.CreatedAt.Month(), comment.CreatedAt.Day()), // mm-dd
	}, nil
}
