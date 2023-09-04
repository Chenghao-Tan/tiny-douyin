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

// 自定义错误类型
var ErrorCommentInaccessible = errors.New("评论不存在或无权访问")

// 评论/删除评论
func Comment(ctx *gin.Context, req *request.CommentReq) (resp *response.CommentResp, err error) {
	// 获取请求用户ID
	req_id, ok := ctx.Get("req_id")
	if !ok {
		utility.Logger().Errorf("ctx.Get (req_id) err: 无法获取")
		return nil, errors.New("无法获取请求用户ID")
	}

	// 存储评论信息
	resp = &response.CommentResp{} // 初始化响应
	if req.Action_Type == 1 {
		// 创建评论
		// 存储评论信息
		comment, err := repo.CreateComment(context.TODO(), req_id.(uint), req.Video_ID, req.Comment_Text)
		if err != nil {
			utility.Logger().Errorf("CreateComment err: %v", err)
			return nil, err
		}

		// 读取评论信息 根据API文档强制要求将其加入响应
		commentInfo, err := readCommentInfo(ctx, comment.ID)
		if err != nil {
			// 响应为评论成功 但评论信息将为空
			utility.Logger().Errorf("readCommentInfo err: %v", err)
		} else {
			// 将该评论加入响应
			resp.Comment = *commentInfo
		}
	} else if req.Action_Type == 2 {
		// 删除评论
		// 删除评论信息
		isReqUsers := repo.CheckUserComments(context.TODO(), req_id.(uint), req.Comment_ID)
		isReqVideos := repo.CheckVideoComments(context.TODO(), req.Video_ID, req.Comment_ID)
		if !isReqUsers || !isReqVideos { // 若非请求用户创建或非处于请求视频下则拒绝删除
			return nil, ErrorCommentInaccessible
		}

		err = repo.DeleteComment(context.TODO(), req.Comment_ID, true) // 永久删除
		if err != nil {
			utility.Logger().Errorf("DeleteComment err: %v", err)
			return nil, err
		}
	} else {
		utility.Logger().Errorf("Invalid action_type err: %v", req.Action_Type)
		return nil, errors.New("操作类型有误")
	}

	return resp, nil
}

// 获取评论列表
func CommentList(ctx *gin.Context, req *request.CommentListReq) (resp *response.CommentListResp, err error) {
	// 读取目标视频评论列表
	comments, err := repo.FindCommentsByCreatedAt(context.TODO(), req.Video_ID, time.Now().Unix(), false, -1) // 倒序向过去查找 不限数量
	if err != nil {
		utility.Logger().Errorf("FindCommentsByCreatedAt err: %v", err)
		return nil, err
	}

	resp = &response.CommentListResp{Comment_List: make([]response.Comment, 0, len(comments))} // 初始化响应
	for _, comment := range comments {
		// 读取评论信息
		commentInfo, err := readCommentInfo(ctx, comment.ID)
		if err != nil {
			utility.Logger().Errorf("readCommentInfo err: %v", err)
			continue // 跳过本条评论
		}

		// 将该评论加入列表
		resp.Comment_List = append(resp.Comment_List, *commentInfo)
	}

	return resp, nil
}
