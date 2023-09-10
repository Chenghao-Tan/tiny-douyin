package repo

import (
	"douyin/repo/internal/db"
	"douyin/repo/internal/db/model"
	"douyin/repo/internal/redis"

	"context"
	"time"
)

// 获取评论主键最大值
func MaxCommentID(ctx context.Context) (id uint, err error) {
	return redis.GetCommentMaxID(ctx)
}

// 创建评论
func CreateComment(ctx context.Context, authorID uint, videoID uint, content string) (comment *model.Comment, err error) {
	comment, err = db.CreateComment(ctx, authorID, videoID, content)
	if err != nil {
		return nil, err
	}
	_ = redis.IncrCommentMaxID(ctx)
	_ = redis.DelUserCommentsCount(ctx, authorID, maxRWTime)
	_ = redis.DelVideoCommentsCount(ctx, videoID, maxRWTime)
	return comment, nil
}

// 删除评论
func DeleteComment(ctx context.Context, id uint, permanently bool) (err error) {
	comment, err2 := ReadCommentBasics(ctx, id) // 读取基本信息以获取作者ID与视频ID (必须在删除前进行)
	err = db.DeleteComment(ctx, id, permanently)
	if err != nil {
		return err
	}
	if err2 == nil { // 若此前成功获取到作者ID与视频ID
		_ = redis.DelUserCommentsCount(ctx, comment.AuthorID, maxRWTime)
		_ = redis.DelVideoCommentsCount(ctx, comment.VideoID, maxRWTime)
	}
	return nil
}

// 读取评论基本信息 (select: ID, CreatedAt, UpdatedAt, Content, AuthorID, VideoID)
func ReadCommentBasics(ctx context.Context, id uint) (comment *model.Comment, err error) {
	comment, err = redis.GetCommentBasics(ctx, id)
	if err == nil { // 命中缓存
		if comment.ID == 0 { // 命中空对象
			time.Sleep(maxRWTime)
			comment, err = redis.GetCommentBasics(ctx, id) // 重试
		} else {
			return comment, nil
		}
	}
	if err == nil { // 命中缓存
		if comment.ID == 0 { // 命中空对象
			return nil, ErrorEmptyObject
		} else {
			return comment, nil
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetCommentBasics(ctx, id, &model.Comment{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadCommentBasics(ctx, id)
		if err == nil {
			_ = redis.SetCommentBasics(ctx, id, record, cacheExpiration)
			return record, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
