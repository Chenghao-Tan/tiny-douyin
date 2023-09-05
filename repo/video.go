package repo

import (
	"douyin/repo/internal/db"
	"douyin/repo/internal/db/model"
	"douyin/repo/internal/redis"

	"context"
	"time"
)

// 获取视频主键最大值
func MaxVideoID(ctx context.Context) (id uint, err error) {
	return redis.GetVideoMaxID(ctx)
}

// 创建视频
func CreateVideo(ctx context.Context, authorID uint, title string) (video *model.Video, err error) {
	video, err = db.CreateVideo(ctx, authorID, title)
	if err != nil {
		return nil, err
	}
	_ = redis.IncrVideoMaxID(ctx)
	_ = redis.DelUserWorksCount(ctx, authorID, maxRWTime)
	return video, nil
}

// 删除视频
func DeleteVideo(ctx context.Context, id uint, permanently bool) (err error) {
	video, err2 := ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID (必须在删除前进行)
	err = db.DeleteVideo(ctx, id, permanently)
	if err != nil {
		return err
	}
	if err2 == nil { // 若此前成功获取到作者ID
		_ = redis.DelUserWorksCount(ctx, video.AuthorID, maxRWTime)
	}
	return nil
}

// 根据创建时间查找视频列表(num==-1时取消数量限制) (select: ID, CreatedAt) //TODO
func FindVideosByCreatedAt(ctx context.Context, createdAt int64, forward bool, num int) (videos []model.Video, err error) {
	return db.FindVideosByCreatedAt(ctx, createdAt, forward, num)
}

// 读取视频基本信息 (select: ID, CreatedAt, UpdatedAt, Title, AuthorID)
func ReadVideoBasics(ctx context.Context, id uint) (video *model.Video, err error) {
	video, err = redis.GetVideoBasics(ctx, id)
	if err == nil { // 命中缓存
		if video.ID == 0 { // 命中空对象
			time.Sleep(maxRWTime)
			video, err = redis.GetVideoBasics(ctx, id) // 重试
		} else {
			return video, nil
		}
	}
	if err == nil { // 命中缓存
		if video.ID == 0 { // 命中空对象
			return nil, ErrorEmptyObject
		} else {
			return video, nil
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetVideoBasics(ctx, id, &model.Video{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadVideoBasics(ctx, id)
		if err == nil {
			_ = redis.SetVideoBasics(ctx, id, record, cacheExpiration)
			return record, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// 读取点赞(用户)列表 (select: Favorited.ID) //TODO
func ReadVideoFavorited(ctx context.Context, id uint) (users []model.User, err error) {
	return db.ReadVideoFavorited(ctx, id)
}

// 读取点赞(用户)数量
func CountVideoFavorited(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetVideoFavoritedCount(ctx, id)
	if err == nil { // 命中缓存
		if count == -1 { // 命中空对象
			time.Sleep(maxRWTime)
			count, err = redis.GetVideoFavoritedCount(ctx, id) // 重试
		} else {
			return count
		}
	}
	if err == nil { // 命中缓存
		if count == -1 { // 命中空对象
			return 0
		} else {
			return count
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetVideoFavoritedCount(ctx, id, -1, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountVideoFavorited(ctx, id)
		if record >= 0 {
			_ = redis.SetVideoFavoritedCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 读取评论列表 (select: Comments.ID) //TODO
func ReadVideoComments(ctx context.Context, id uint) (comments []model.Comment, err error) {
	return db.ReadVideoComments(ctx, id)
}

// 读取评论数量
func CountVideoComments(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetVideoCommentsCount(ctx, id)
	if err == nil { // 命中缓存
		if count == -1 { // 命中空对象
			time.Sleep(maxRWTime)
			count, err = redis.GetVideoCommentsCount(ctx, id) // 重试
		} else {
			return count
		}
	}
	if err == nil { // 命中缓存
		if count == -1 { // 命中空对象
			return 0
		} else {
			return count
		}
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetVideoCommentsCount(ctx, id, -1, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountVideoComments(ctx, id)
		if record >= 0 {
			_ = redis.SetVideoCommentsCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 检查评论所属 //TODO
func CheckVideoComments(ctx context.Context, id uint, commentID uint) (isIts bool) {
	return db.CheckVideoComments(ctx, id, commentID)
}
