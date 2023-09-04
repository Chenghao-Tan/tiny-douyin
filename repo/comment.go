package repo

import (
	"douyin/repo/db"
	"douyin/repo/db/model"
	"douyin/repo/redis"

	"context"
)

// 创建评论
func CreateComment(ctx context.Context, authorID uint, videoID uint, content string) (comment *model.Comment, err error) {
	comment, err = db.CreateComment(ctx, authorID, videoID, content)
	if err != nil {
		return nil, err
	}
	_ = redis.DelUserCommentsCount(ctx, authorID, maxWriteTime)
	_ = redis.DelVideoCommentsCount(ctx, videoID, maxWriteTime)
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
		_ = redis.DelUserCommentsCount(ctx, comment.AuthorID, maxWriteTime)
		_ = redis.DelVideoCommentsCount(ctx, comment.VideoID, maxWriteTime)
	}
	return nil
}

// 根据视频ID和创建时间查找评论列表(num==-1时取消数量限制) (select: ID, CreatedAt) //TODO
func FindCommentsByCreatedAt(ctx context.Context, videoID uint, createdAt int64, forward bool, num int) (comments []model.Comment, err error) {
	return db.FindCommentsByCreatedAt(ctx, videoID, createdAt, forward, num)
}

// 读取评论基本信息 (select: *)
func ReadCommentBasics(ctx context.Context, id uint) (comment *model.Comment, err error) {
	comment, err = redis.GetCommentBasics(ctx, id)
	if err == nil { // 命中缓存
		return comment, nil
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetCommentBasics(ctx, id, &model.Comment{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadCommentBasics(ctx, id)
		if err == nil {
			_ = redis.SetCommentBasics(ctx, id, record, cacheExpiration)
			return record, nil
		}
	}
	return nil, err // 当出现错误
}
