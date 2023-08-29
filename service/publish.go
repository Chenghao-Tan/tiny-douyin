package service

import (
	"douyin/repo/db"
	"douyin/repo/oss"
	"douyin/service/type/request"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 发布视频
func Publish(ctx *gin.Context, req *request.PublishReq) (resp *response.PublishResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("req_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (req_id) rr: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 先尝试打开文件 若无法打开则不创建数据库条目
	videoStream, err := req.Data.Open()
	if err != nil {
		utility.Logger().Errorf("file.Open err: %v", err)
		return nil, err
	}
	defer videoStream.Close() // 不保证自动关闭成功

	// 存储视频信息
	video, err := db.CreateVideo(context.TODO(), req_id.(uint), req.Title)
	if err != nil {
		utility.Logger().Errorf("CreateVideo err: %v", err)
		return nil, err
	}

	// 上传视频数据(封面为默认)
	err = oss.UploadVideoStream(context.TODO(), strconv.FormatUint(uint64(video.ID), 10), videoStream, req.Data.Size)
	if err != nil {
		utility.Logger().Errorf("UploadVideoStream err: %v", err)
		return nil, err
	}

	// 创建更新封面异步任务
	go func() {
		err2 := oss.UpdateCover(context.TODO(), strconv.FormatUint(uint64(video.ID), 10)) // 不保证自动更新成功
		if err2 != nil {
			utility.Logger().Errorf("UpdateCover err: %v", err2)
		}
	}()

	return &response.PublishResp{}, nil
}

// 获取发布列表
func PublishList(ctx *gin.Context, req *request.PublishListReq) (resp *response.PublishListResp, err error) {
	// 读取目标用户信息
	works, err := db.ReadUserWorks(context.TODO(), req.User_ID)
	if err != nil {
		utility.Logger().Errorf("ReadUserWorks err: %v", err)
		return nil, err
	}

	// 读取目标用户发布列表
	resp = &response.PublishListResp{Video_List: make([]response.Video, 0, len(works))} // 初始化响应
	for _, video := range works {
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
