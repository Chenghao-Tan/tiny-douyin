package repo

import (
	"douyin/repo/db"
	"douyin/repo/db/model"
	"douyin/repo/redis"

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

// 读取用户基本信息 (select: *)
func ReadUserBasics(ctx context.Context, id uint) (user *model.User, err error) {
	user, err = redis.GetUserBasics(ctx, id)
	if err == nil { // 命中缓存
		return user, nil
	}
	if err == redis.ErrorRedisNil { // 启动同步
		_ = redis.SetUserBasics(ctx, id, &model.User{}, emptyExpiration) // 防止缓存穿透与缓存击穿
		record, err := db.ReadUserBasics(ctx, id)
		if err != nil {
			_ = redis.SetUserBasics(ctx, id, record, cacheExpiration)
			return record, nil
		}
	}
	return nil, err // 当出现错误
}

// 创建点赞关系
func CreateUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("fav:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(videoID), 10) + ":1")

	video, err := db.ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID
	if err != nil {
		return err
	}
	return redis.SetUserFavorites(ctx, id, videoID, video.AuthorID, true, maxSyncDelay)
}

// 删除点赞关系
func DeleteUserFavorites(ctx context.Context, id uint, videoID uint) (err error) {
	// 加入同步队列
	syncQueue.Push("fav:" + strconv.FormatUint(uint64(id), 10) + ":" + strconv.FormatUint(uint64(videoID), 10) + ":0")

	video, err := db.ReadVideoBasics(ctx, id) // 读取基本信息以获取作者ID
	if err != nil {
		return err
	}
	return redis.SetUserFavorites(ctx, id, videoID, video.AuthorID, false, maxSyncDelay)
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
		}
	}
	return 0 // 当出现错误
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
		}
	}
	return 0 // 当出现错误
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
	}
	return false // 当出现错误
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
		}
	}
	return 0 // 当出现错误
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
		}
	}
	return 0 // 当出现错误
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
	}
	return false // 当出现错误
}
