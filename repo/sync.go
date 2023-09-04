package repo

import (
	"douyin/repo/internal/db"
	"douyin/utility"

	"container/list"
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

type MessageQueue struct {
	queue *list.List
	lock  *sync.Mutex
}

func (mq *MessageQueue) Init() {
	mq.queue = list.New()
	mq.lock = &sync.Mutex{}
}

func (mq *MessageQueue) Push(message any) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.queue.PushBack(message)
}

func (mq *MessageQueue) Pop() (message any) {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if mq.queue.Front() != nil {
		return mq.queue.Remove(mq.queue.Front())
	}
	return nil
}

func (mq *MessageQueue) Len() (len int) {
	return mq.queue.Len()
}

var syncQueue = &MessageQueue{} // 基于链表的简易消息队列
var syncCron = cron.New()       // 定时任务规划器

func syncTask() {
	// 取出此刻所有消息
	opsNum := syncQueue.Len()
	successCount := 0
	for i := 0; i < opsNum; i++ {
		message := syncQueue.Pop()
		if message != nil {
			op, ok := message.(string)
			if !ok {
				continue
			}
			split := strings.Split(op, ":")
			if len(split) != 4 {
				continue
			}

			if split[0] == "fav" { // 同步点赞变更
				userID, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为用户ID", split[1])
					continue
				}
				videoID, err := strconv.ParseUint(split[2], 10, 64)
				if err != nil {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为视频ID", split[2])
					continue
				}
				isFavorite := split[3]

				if isFavorite == "1" {
					err := db.CreateUserFavorites(context.TODO(), uint(userID), uint(videoID))
					if err != nil {
						utility.Logger().Errorf("repo.syncTask (CreateUserFavorites) err: %v", err)
					} else {
						successCount++
					}
				} else if isFavorite == "0" {
					err := db.DeleteUserFavorites(context.TODO(), uint(userID), uint(videoID))
					if err != nil {
						utility.Logger().Errorf("repo.syncTask (DeleteUserFavorites) err: %v", err)
					} else {
						successCount++
					}
				} else {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为点赞信息", isFavorite)
				}
			}

			if split[0] == "flw" { // 同步关注变更
				userID, err := strconv.ParseUint(split[1], 10, 64)
				if err != nil {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为用户ID", split[1])
					continue
				}
				followID, err := strconv.ParseUint(split[2], 10, 64)
				if err != nil {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为被关注用户ID", split[2])
					continue
				}
				isFollowing := split[3]

				if isFollowing == "1" {
					err := db.CreateUserFollows(context.TODO(), uint(userID), uint(followID))
					if err != nil {
						utility.Logger().Errorf("repo.syncTask (CreateUserFollows) err: %v", err)
					} else {
						successCount++
					}
				} else if isFollowing == "0" {
					err := db.DeleteUserFollows(context.TODO(), uint(userID), uint(followID))
					if err != nil {
						utility.Logger().Errorf("repo.syncTask (DeleteUserFollows) err: %v", err)
					} else {
						successCount++
					}
				} else {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为关注信息", isFollowing)
				}
			}
		}
	}

	if opsNum > 0 {
		utility.Logger().Infof("repo.syncTask info: %v项同步成功", successCount)
	}
}
