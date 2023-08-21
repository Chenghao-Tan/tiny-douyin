package service

import (
	"douyin/repo/db"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 视频流
func Feed(ctx *gin.Context, req *request.FeedReq) (resp *response.FeedResp, err error) {
	// 读取视频列表 //TODO
	latest_time, err := strconv.ParseInt(req.Latest_Time, 10, 64)
	if err != nil {
		utility.Logger().Errorf("ParseInt err: %v", err)
		return nil, err
	}
	videos, err := db.FindVideosByUpdatedAt(context.TODO(), latest_time, false, 30) // 倒序向过去查找 最多30条
	if err != nil {
		utility.Logger().Errorf("FindVideosByUpdatedAt err: %v", err)
		return nil, err
	}

	// 初始化响应
	resp = &response.FeedResp{
		// Next_Time: 0, // 本次返回的视频中发布最早的时间 根据API文档默认为不发送
	}
	if len(videos) > 0 { // 如果查找结果中有视频
		resp.Next_Time = uint(videos[len(videos)-1].UpdatedAt.Unix() * 1000) // 更新该时间戳 API文档有误 响应实为毫秒时间戳 故在此转换
	}

	// 向响应中添加视频
	for _, video := range videos {
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
