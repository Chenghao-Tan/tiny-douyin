package repo

import (
	"douyin/repo/internal/db"
	"douyin/repo/internal/redis"
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

func syncTask() { // 回写策略的缓存持久化/一致性同步任务
	opsNum := syncQueue.Len() // 总项目数量(此刻队列内消息数)

	// 取消息
	type relation struct {
		id1 uint
		id2 uint
	}
	favHash := make(map[relation]bool)
	flwHash := make(map[relation]bool)
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
					favHash[relation{uint(userID), uint(videoID)}] = true
				} else if isFavorite == "0" {
					favHash[relation{uint(userID), uint(videoID)}] = false
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
					flwHash[relation{uint(userID), uint(followID)}] = true
				} else if isFollowing == "0" {
					flwHash[relation{uint(userID), uint(followID)}] = false
				} else {
					utility.Logger().Errorf("repo.syncTask err: %v无法识别为关注信息", isFollowing)
				}
			}
		}
	}

	// 持久化
	opsNum = len(favHash) + len(flwHash) // 总项目数量(仅有效)
	successCount := 0
	favCreateIDs := make([][]uint, 2)
	favDeleteIDs := make([][]uint, 2)
	for key, value := range favHash {
		if value {
			favCreateIDs[0] = append(favCreateIDs[0], key.id1)
			favCreateIDs[1] = append(favCreateIDs[1], key.id2)
		} else {
			favDeleteIDs[0] = append(favDeleteIDs[0], key.id1)
			favDeleteIDs[1] = append(favDeleteIDs[1], key.id2)
		}
	}
	flwCreateIDs := make([][]uint, 2)
	flwDeleteIDs := make([][]uint, 2)
	for key, value := range flwHash {
		if value {
			flwCreateIDs[0] = append(flwCreateIDs[0], key.id1)
			flwCreateIDs[1] = append(flwCreateIDs[1], key.id2)
		} else {
			flwDeleteIDs[0] = append(flwDeleteIDs[0], key.id1)
			flwDeleteIDs[1] = append(flwDeleteIDs[1], key.id2)
		}
	}
	successCount += int(db.CreateUserFavoritesBatch(context.TODO(), favCreateIDs[0], favCreateIDs[1]))
	successCount += int(db.DeleteUserFavoritesBatch(context.TODO(), favDeleteIDs[0], favDeleteIDs[1]))
	successCount += int(db.CreateUserFollowsBatch(context.TODO(), flwCreateIDs[0], flwCreateIDs[1]))
	successCount += int(db.DeleteUserFollowsBatch(context.TODO(), flwDeleteIDs[0], flwDeleteIDs[1]))

	if successCount < opsNum { // 无任务时successCount==opsNum==0
		utility.Logger().Errorf("repo.syncTask err: %v项同步失败", opsNum-successCount)
	}
}

func updateTask() { // 直写策略的(特殊)缓存一致性同步任务
	opsNum := 4 // 总项目数量

	// 各模型主键最大值同步
	successCount := 0
	userMaxID, err := db.MaxUserID(context.TODO())
	if err == nil {
		if redis.SetUserMaxID(context.TODO(), userMaxID) == nil {
			successCount++
		}
	}
	videoMaxID, err := db.MaxVideoID(context.TODO())
	if err == nil {
		if redis.SetVideoMaxID(context.TODO(), videoMaxID) == nil {
			successCount++
		}
	}
	commentMaxID, err := db.MaxCommentID(context.TODO())
	if err == nil {
		if redis.SetCommentMaxID(context.TODO(), commentMaxID) == nil {
			successCount++
		}
	}
	messageMaxID, err := db.MaxMessageID(context.TODO()) // 由于暂未选定消息缓存策略, 暂时在此同步 //TODO
	if err == nil {
		if redis.SetMessageMaxID(context.TODO(), messageMaxID) == nil {
			successCount++
		}
	}

	if successCount < opsNum {
		utility.Logger().Errorf("repo.updateTask err: %v项同步失败", opsNum-successCount)
	}
}
