package redis

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const prefixUserFavorites = "user:fav:"                            // 后接三十六进制userID (节约key长度)
const prefixUserFavoritesDelta = prefixUserFavorites + "delta:"    // 后接三十六进制userID:videoID (节约key长度)
const prefixUserFavoritesCount = prefixUserFavorites + "count:"    // 后接三十六进制userID (节约key长度)
const prefixUserFavoritedCount = prefixUserFavorites + "dcount:"   // 后接三十六进制userID (节约key长度)
const prefixVideoFavorited = "video:fav:"                          // 暂只用于构建其他前缀
const prefixVideoFavoritedCount = prefixVideoFavorited + "dcount:" // 后接三十六进制videoID (节约key长度)

// 设置点赞关系变更记录(并设置相关计数)
func setUserFavoritesDelta(ctx context.Context, userID uint, videoID uint, authorID uint, isFavorite bool, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		deltaKey := prefixUserFavoritesDelta + strconv.FormatUint(uint64(userID), 36) + ":" + strconv.FormatUint(uint64(videoID), 36)
		countKey := prefixUserFavoritesCount + strconv.FormatUint(uint64(userID), 36)
		dcountKey := prefixUserFavoritedCount + strconv.FormatUint(uint64(authorID), 36)
		videoDcountKey := prefixVideoFavoritedCount + strconv.FormatUint(uint64(videoID), 36)

		if isFavorite {
			pipe.SetEx(ctx, deltaKey, isFavorite, expiration) // 在确定数据库已被写入后过期
			pipe.Incr(ctx, countKey)
			pipe.Incr(ctx, dcountKey)
			pipe.Incr(ctx, videoDcountKey)
		} else {
			pipe.SetEx(ctx, deltaKey, isFavorite, expiration) // 在确定数据库已被写入后过期
			pipe.Decr(ctx, countKey)
			pipe.Decr(ctx, dcountKey)
			pipe.Decr(ctx, videoDcountKey)
		}
		pipe.Expire(ctx, countKey, expiration)       // 在确定数据库已被写入后(即最后一条变更记录过期后)强制刷新
		pipe.Expire(ctx, dcountKey, expiration)      // 在确定数据库已被写入后(即最后一条变更记录过期后)强制刷新
		pipe.Expire(ctx, videoDcountKey, expiration) // 在确定数据库已被写入后(即最后一条变更记录过期后)强制刷新

		return nil
	})
	return err
}

// 读取点赞关系变更记录
func getUserFavoritesDelta(ctx context.Context, userID uint, videoID uint) (isFavorite bool, err error) {
	key := prefixUserFavoritesDelta + strconv.FormatUint(uint64(userID), 36) + ":" + strconv.FormatUint(uint64(videoID), 36)
	return _redis.Get(ctx, key).Bool()
}

// 设置点赞关系(仅用于一致性同步时修正主记录)
func SetUserFavoritesBit(ctx context.Context, userID uint, videoID uint, isFavorite bool) (err error) {
	key := prefixUserFavorites + strconv.FormatUint(uint64(userID), 36)
	value := 0
	if isFavorite {
		value = 1
	}
	return _redis.SetBit(ctx, key, int64(videoID), value).Err()
}

// 设置点赞关系(仅用于处理用户请求 会导致随机不信任缓存暂时禁用)
func SetUserFavorites(ctx context.Context, userID uint, videoID uint, authorID uint, isFavorite bool, maxSyncDelay time.Duration) (err error) {
	key := prefixUserFavoritesDelta + strconv.FormatUint(uint64(userID), 36) + ":" + strconv.FormatUint(uint64(videoID), 36)
	value, err := _redis.Get(ctx, key).Bool() // 读取变更记录以过滤重复请求
	if err != nil && err != ErrorRedisNil {
		return err
	}
	if err != ErrorRedisNil && value && isFavorite { // 已设置过相同变更
		return ErrorRecordExists // 防止重复计数
	}
	if err != ErrorRedisNil && !value && !isFavorite { // 已设置过相同变更
		return ErrorRecordNotExists // 防止重复计数
	}

	// 写入变更记录 在最长写入数据库用时+1秒时过期以确保数据库已写入 过期前禁用随机不信任缓存以防错误同步
	err = setUserFavoritesDelta(ctx, userID, videoID, authorID, isFavorite, maxSyncDelay+time.Second)
	if err != nil {
		// 一般为事务整体失败
		return err
	}

	// 主记录在变更记录过期前100毫秒时写入以应对可能即将到来的访问 并覆盖此前的所有错误写入 原因参考缓存双删
	go func() {
		time.Sleep(maxSyncDelay + time.Millisecond*900)

		_ = SetUserFavoritesBit(ctx, userID, videoID, isFavorite)
	}()

	return nil
}

// 读取点赞关系
func GetUserFavorites(ctx context.Context, userID uint, videoID uint, distrustProbability float32) (isFavorite bool, err error) {
	isFavorite, err = getUserFavoritesDelta(ctx, userID, videoID)
	if err == nil { // 若有变更记录存在则直接返回(此时禁用随机不信任缓存)
		return isFavorite, nil
	}

	switch {
	case distrustProbability == 0:
		// 强制信任缓存
	case (distrustProbability > 0 && distrustProbability < 1):
		// 随机不信任缓存
		if rand.Intn(int(1/distrustProbability)) == 0 {
			return false, ErrorRedisNil // 返回查找结果为空, 以供触发一致性同步
		}
	case distrustProbability == 1:
		// 强制不信任缓存
		return false, ErrorRedisNil
	default:
		return false, errors.New("distrustProbability必须在0-1之间")
	}

	// 从主记录正常读取
	key := prefixUserFavorites + strconv.FormatUint(uint64(userID), 36)
	value, err := _redis.GetBit(ctx, key, int64(videoID)).Result()
	if err != nil {
		return false, err
	}
	if value == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// 设置用户点赞数
func SetUserFavoritesCount(ctx context.Context, userID uint, count int64, expiration time.Duration) (err error) {
	key := prefixUserFavoritesCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取用户点赞数
func GetUserFavoritesCount(ctx context.Context, userID uint) (count int64, err error) {
	key := prefixUserFavoritesCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Int64()
}

// 设置用户受赞数
func SetUserFavoritedCount(ctx context.Context, userID uint, count int64, expiration time.Duration) (err error) {
	key := prefixUserFavoritedCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取用户受赞数
func GetUserFavoritedCount(ctx context.Context, userID uint) (count int64, err error) {
	key := prefixUserFavoritedCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Int64()
}

// 设置视频受赞数
func SetVideoFavoritedCount(ctx context.Context, userID uint, count int64, expiration time.Duration) (err error) {
	key := prefixVideoFavoritedCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取视频受赞数
func GetVideoFavoritedCount(ctx context.Context, userID uint) (count int64, err error) {
	key := prefixVideoFavoritedCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Int64()
}
