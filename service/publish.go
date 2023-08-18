package service

import (
	"douyin/repository/dao"
	"douyin/repository/model"
	"douyin/service/types/request"
	"douyin/service/types/response"
	"douyin/utils"
	"douyin/utils/oss"

	"context"
	"errors"
	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 发布视频
func Publish(ctx *gin.Context, req *request.PublishReq, file *multipart.FileHeader) (resp *response.PublishResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("user_id")
	if !ok {
		utils.ZapLogger.Errorf("ctx.Get (user_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 先尝试打开文件 若无法打开则不创建数据库条目
	videoStream, err := file.Open()
	if err != nil {
		utils.ZapLogger.Errorf("file.Open err: %v", err)
		return nil, err
	}
	defer videoStream.Close() // 不保证自动关闭成功

	// 准备要存储的内容
	video := &model.Video{
		Title:  req.Title,
		UserID: req_id.(uint),
	}

	// 存储视频信息 //TODO
	video, err = dao.CreateVideo(context.TODO(), video)
	if err != nil {
		utils.ZapLogger.Errorf("CreateVideo err: %v", err)
		return nil, err
	}

	// 上传视频数据(封面为默认)
	err = oss.UploadVideoStream(context.TODO(), strconv.FormatUint(uint64(video.ID), 10), videoStream, file.Size)
	if err != nil {
		utils.ZapLogger.Errorf("UploadVideoStream err: %v", err)
		return nil, err
	}

	// 创建更新封面异步任务
	go func() {
		oss.UpdateCover(context.TODO(), strconv.FormatUint(uint64(video.ID), 10)) // 不保证自动更新成功
	}()

	return &response.PublishResp{}, nil
}

// 获取发布列表
func PublishList(ctx *gin.Context, req *request.PublishListReq) (resp *response.PublishListResp, err error) {
	// 读取目标用户信息
	user_id, err := strconv.ParseUint(req.User_ID, 10, 64)
	if err != nil {
		utils.ZapLogger.Errorf("ParseUint err: %v", err)
		return nil, err
	}
	user, err := dao.FindUserByID(context.TODO(), uint(user_id))
	if err != nil {
		utils.ZapLogger.Errorf("FindUserByID err: %v", err)
		return nil, err
	}

	// 读取目标用户发布列表 //TODO
	resp = &response.PublishListResp{} // 初始化响应
	for _, video := range user.Works {
		// 读取视频信息
		videoInfo, err := readVideoInfo(ctx, video.ID)
		if err != nil {
			utils.ZapLogger.Errorf("readVideoInfo err: %v", err)
			continue // 跳过本条视频
		}

		// 将该视频加入列表
		resp.Video_List = append(resp.Video_List, *videoInfo)
	}

	return resp, nil
}
