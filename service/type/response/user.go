package response

type UserRegisterResp struct {
	Status
	User_ID uint   `json:"user_id"` // 用户id
	Token   string `json:"token"`   // 用户鉴权token
}

type UserLoginResp struct {
	Status
	User_ID uint   `json:"user_id"` // 用户id
	Token   string `json:"token"`   // 用户鉴权token
}

type UserInfoResp struct {
	Status
	User User `json:"user"` // 用户信息
}
