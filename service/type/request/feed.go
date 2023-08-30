package request

type FeedReq struct {
	Latest_Time int64  `json:"latest_time" form:"latest_time" binding:"omitempty,min=0"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	Token       string `json:"token" form:"token" binding:"omitempty,jwt"`               // 可选参数，用户登录状态下设置
}
