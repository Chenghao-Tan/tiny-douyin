package repo

import (
	"douyin/repo/internal/db"
	"douyin/repo/internal/db/model"
	"douyin/repo/internal/redis"

	"context"
	"strconv"
)

// 创建用户
func CreateUser(ctx context.Context, username string, password string, signature string) (user *model.User, err error) {
	return db.CreateUser(ctx, username, password, signature)
}

// 检查用户名是否可用
func CheckUserRegister(ctx context.Context, username string) (isAvailable bool) {
	return db.CheckUserRegister(ctx, username)
}

// 检查用户名和密码是否有效
func CheckUserLogin(ctx context.Context, username string, password string) (id uint, isValid bool) {
	return db.CheckUserLogin(ctx, username, password)
}

// 读取用户基本信息 (select: ID, CreatedAt, UpdatedAt, Username, Signature)
func ReadUserBasics(ctx context.Context, id uint) (user *model.User, err error) {
	user, err = redis.GetUserBasics(ctx, id)
	if err == nil { // 命中缓存
		return user, nil
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserBasics(ctx, id, &model.User{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadUserBasics(ctx, id)
		if err == nil {
			_ = redis.SetUserBasics(ctx, id, record, cacheExpiration)
			return record, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// 读取作品(视频)列表 (select: Works.ID) //TODO
func ReadUserWorks(ctx context.Context, id uint) (videos []model.Video, err error) {
	return db.ReadUserWorks(ctx, id)
}

// 读取作品(视频)数量
func CountUserWorks(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserWorksCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserWorksCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserWorks(ctx, id)
		if record >= 0 {
			_ = redis.SetUserWorksCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 创建点赞关系
func CreateUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("fav:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(videoID), 10) + ":1")

	video, err := ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID
	if err != nil {
		return err
	}
	return redis.SetUserFavorites(ctx, id, videoID, video.AuthorID, true, maxSyncDelay)
}

// 删除点赞关系
func DeleteUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("fav:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(videoID), 10) + ":0")

	video, err := ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID
	if err != nil {
		return err
	}
	return redis.SetUserFavorites(ctx, id, videoID, video.AuthorID, false, maxSyncDelay)
}

// 读取点赞(视频)列表 (select: Favorites.ID) //TODO
func ReadUserFavorites(ctx context.Context, id uint) (videos []model.Video, err error) {
	return db.ReadUserFavorites(ctx, id)
}

// 读取点赞(视频)数量
func CountUserFavorites(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserFavoritesCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserFavoritesCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserFavorites(ctx, id)
		if record >= 0 {
			_ = redis.SetUserFavoritesCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 读取获赞数量
func CountUserFavorited(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserFavoritedCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserFavoritedCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserFavorited(ctx, id)
		if record >= 0 {
			_ = redis.SetUserFavoritedCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 检查点赞关系
func CheckUserFavorites(ctx context.Context, id uint, videoID uint) (isFavorite bool) {
	isFavorite, err := redis.GetUserFavorites(ctx, id, videoID, distrustProbability)
	if err == nil { // 命中缓存
		return isFavorite
	}
	if err == redis.ErrorRedisNil { // 启动同步
		record := db.CheckUserFavorites(ctx, id, videoID)
		_ = redis.SetUserFavoritesBit(ctx, id, videoID, record) // 立即修正缓存主记录
		return record
	} else {
		return false
	}
}

// 读取评论列表 (select: Comments.ID) //TODO
func ReadUserComments(ctx context.Context, id uint) (comments []model.Comment, err error) {
	return db.ReadUserComments(ctx, id)
}

// 读取评论数量
func CountUserComments(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserCommentsCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserCommentsCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserComments(ctx, id)
		if record >= 0 {
			_ = redis.SetUserCommentsCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 检查评论所属 //TODO
func CheckUserComments(ctx context.Context, id uint, commentID uint) (isIts bool) {
	return db.CheckUserComments(ctx, id, commentID)
}

// 创建关注关系
func CreateUserFollows(ctx context.Context, id uint, followID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("flw:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(followID), 10) + ":1")

	return redis.SetUserFollows(ctx, id, followID, true, maxSyncDelay)
}

// 删除关注关系
func DeleteUserFollows(ctx context.Context, id uint, followID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("flw:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(followID), 10) + ":0")

	return redis.SetUserFollows(ctx, id, followID, false, maxSyncDelay)
}

// 读取关注(用户)列表 (select: Follows.ID) //TODO
func ReadUserFollows(ctx context.Context, id uint) (users []model.User, err error) {
	return db.ReadUserFollows(ctx, id)
}

// 读取关注(用户)数量
func CountUserFollows(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserFollowsCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserFollowsCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserFollows(ctx, id)
		if record >= 0 {
			_ = redis.SetUserFollowsCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 读取粉丝(用户)列表 (select: Followers.ID) //TODO
func ReadUserFollowers(ctx context.Context, id uint) (users []model.User, err error) {
	return db.ReadUserFollowers(ctx, id)
}

// 读取粉丝(用户)数量
func CountUserFollowers(ctx context.Context, id uint) (count int64) {
	count, err := redis.GetUserFollowersCount(ctx, id)
	if err == nil { // 命中缓存
		return count
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserFollowersCount(ctx, id, 0, emptyExpiration) // 防止缓存穿透与缓存击穿
		record := db.CountUserFollowers(ctx, id)
		if record >= 0 {
			_ = redis.SetUserFollowersCount(ctx, id, record, cacheExpiration)
			return record
		} else {
			return -1
		}
	} else {
		return -1
	}
}

// 检查关注关系
func CheckUserFollows(ctx context.Context, id uint, followID uint) (isFollowing bool) {
	isFavorite, err := redis.GetUserFollows(ctx, id, followID, distrustProbability)
	if err == nil { // 命中缓存
		return isFavorite
	}
	if err == redis.ErrorRedisNil { // 启动同步
		record := db.CheckUserFollows(ctx, id, followID)
		_ = redis.SetUserFollowsBit(ctx, id, followID, record) // 立即修正缓存主记录
		return record
	} else {
		return false
	}
}

// 读取消息列表 (select: Messages.ID) //TODO
func ReadUserMessages(ctx context.Context, id uint) (messages []model.Message, err error) {
	return db.ReadUserMessages(ctx, id)
}

// 计算消息数量 //TODO
func CountUserMessages(ctx context.Context, id uint) (count int64) {
	return db.CountUserMessages(ctx, id)
}
