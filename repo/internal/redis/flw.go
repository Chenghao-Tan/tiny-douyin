package redis

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const prefixUserFollows = "user:flw:"                          // 后接三十六进制userID (节约key长度)
const prefixUserFollowsDelta = prefixUserFollows + "delta:"    // 后接三十六进制userID:followID (节约key长度)
const prefixUserFollowsCount = prefixUserFollows + "count:"    // 后接三十六进制userID (节约key长度)
const prefixUserFollowersCount = prefixUserFollows + "dcount:" // 后接三十六进制followID (节约key长度)

// 设置关注关系变更记录(并设置相关计数)
func setUserFollowsDelta(ctx context.Context, userID uint, followID uint, isFollowing bool, expiration time.Duration) (err error) {
	_, err = _redis.TxPipelined(ctx, func(pipe redis.Pipeliner) error { // 使用事务
		deltaKey := prefixUserFollowsDelta + strconv.FormatUint(uint64(userID), 36) + ":" + strconv.FormatUint(uint64(followID), 36)
		countKey := prefixUserFollowsCount + strconv.FormatUint(uint64(userID), 36)
		dcountKey := prefixUserFollowersCount + strconv.FormatUint(uint64(followID), 36)

		if isFollowing {
			pipe.SetEx(ctx, deltaKey, isFollowing, expiration) // 在确定数据库已被写入后过期
			pipe.Incr(ctx, countKey)
			pipe.Incr(ctx, dcountKey)
		} else {
			pipe.SetEx(ctx, deltaKey, isFollowing, expiration) // 在确定数据库已被写入后过期
			pipe.Decr(ctx, countKey)
			pipe.Decr(ctx, dcountKey)
		}
		pipe.Expire(ctx, countKey, expiration)  // 在确定数据库已被写入后(即最后一条变更记录过期后)强制刷新
		pipe.Expire(ctx, dcountKey, expiration) // 在确定数据库已被写入后(即最后一条变更记录过期后)强制刷新

		return nil
	})
	return err
}

// 读取关注关系变更记录
func getUserFollowsDelta(ctx context.Context, userID uint, followID uint) (isFollowing bool, err error) {
	key := prefixUserFollowsDelta + strconv.FormatUint(uint64(userID), 36) + ":" + strconv.FormatUint(uint64(followID), 36)
	return _redis.Get(ctx, key).Bool()
}

// 设置关注关系(仅用于一致性同步时修正主记录)
func SetUserFollowsBit(ctx context.Context, userID uint, followID uint, isFollowing bool) (err error) {
	if userID == followID {
		return ErrorSelfFollow // 默认禁止自己关注自己
	}

	key := prefixUserFollows + strconv.FormatUint(uint64(userID), 36)
	value := 0
	if isFollowing {
		value = 1
	}
	return _redis.SetBit(ctx, key, int64(followID), value).Err()
}

// 设置关注关系(仅用于处理用户请求 会导致随机不信任缓存暂时禁用)
func SetUserFollows(ctx context.Context, userID uint, followID uint, isFollowing bool, maxSyncDelay time.Duration) (err error) {
	if userID == followID {
		return ErrorSelfFollow // 默认禁止自己关注自己
	}

	key := prefixUserFollowsDelta + strconv.FormatUint(uint64(userID), 36) + ":" + strconv.FormatUint(uint64(followID), 36)
	value, err := _redis.Get(ctx, key).Bool() // 读取变更记录以过滤重复请求
	if err != nil && err != ErrorRedisNil {
		return err
	}
	if err != ErrorRedisNil && value && isFollowing { // 已设置过相同变更
		return ErrorRecordExists // 防止重复计数
	}
	if err != ErrorRedisNil && !value && !isFollowing { // 已设置过相同变更
		return ErrorRecordNotExists // 防止重复计数
	}

	// 写入变更记录 在最长同步延迟+1秒时过期以确保缓存已写入(防止因精确到秒向下取整导致的问题) 过期前禁用随机不信任缓存以防错误同步
	err = setUserFollowsDelta(ctx, userID, followID, isFollowing, maxSyncDelay+time.Second)
	if err != nil {
		// 一般为事务整体失败
		return err
	}

	// 主记录在最大同步延迟后写入以应对可能即将到来的访问 并覆盖此前的所有错误写入 原因参考缓存双删
	go func() {
		time.Sleep(maxSyncDelay)

		_ = SetUserFollowsBit(ctx, userID, followID, isFollowing)
	}()

	return nil
}

// 读取关注关系
func GetUserFollows(ctx context.Context, userID uint, followID uint, distrustProbability float32) (isFollowing bool, err error) {
	if userID == followID {
		return false, nil // 默认自己不关注自己
	}

	isFollowing, err = getUserFollowsDelta(ctx, userID, followID)
	if err == nil { // 若有变更记录存在则直接返回(此时禁用随机不信任缓存)
		return isFollowing, nil
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
	key := prefixUserFollows + strconv.FormatUint(uint64(userID), 36)
	value, err := _redis.GetBit(ctx, key, int64(followID)).Result()
	if err != nil {
		return false, err
	}
	if value == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

// 设置用户关注数
func SetUserFollowsCount(ctx context.Context, userID uint, count int64, expiration time.Duration) (err error) {
	key := prefixUserFollowsCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取用户关注数
func GetUserFollowsCount(ctx context.Context, userID uint) (count int64, err error) {
	key := prefixUserFollowsCount + strconv.FormatUint(uint64(userID), 36)
	return _redis.Get(ctx, key).Int64()
}

// 设置用户粉丝数
func SetUserFollowersCount(ctx context.Context, followID uint, count int64, expiration time.Duration) (err error) {
	key := prefixUserFollowersCount + strconv.FormatUint(uint64(followID), 36)
	return _redis.SetEx(ctx, key, count, randomExpiration(expiration)).Err()
}

// 读取用户粉丝数
func GetUserFollowersCount(ctx context.Context, followID uint) (count int64, err error) {
	key := prefixUserFollowersCount + strconv.FormatUint(uint64(followID), 36)
	return _redis.Get(ctx, key).Int64()
}
