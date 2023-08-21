package response

type FeedResp struct {
	Status
	Next_Time  uint    `json:"next_time"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time API文档有误 实为毫秒时间戳
	Video_List []Video `json:"video_list"` // 视频列表
}
