package response

type PublishResp struct {
	Status
}

type PublishListResp struct {
	Status
	Video_List []Video `json:"video_list"` // 用户发布的视频列表
}
