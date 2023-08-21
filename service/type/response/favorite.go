package response

type FavoriteResp struct {
	Status
}

type FavoriteListResp struct {
	Status
	Video_List []Video `json:"video_list"` // 用户点赞视频列表
}
