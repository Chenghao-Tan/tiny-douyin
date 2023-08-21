package response

type FollowResp struct {
	Status
}

type FollowListResp struct {
	Status
	User_List []User `json:"user_list"` // 用户信息列表
}

type FollowerListResp struct {
	Status
	User_List []User `json:"user_list"` // 用户信息列表
}

type FriendListResp struct {
	Status
	User_List []FriendUser `json:"user_list"` // 用户(好友)信息列表
}
