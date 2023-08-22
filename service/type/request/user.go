package request

type UserRegisterReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"`       // 注册用户名 1-32个字符
	Password string `json:"password" form:"password" binding:"required,min=6,max=32"` // 注册密码 6-32个字符
}

type UserLoginReq struct {
	Username string `json:"username" form:"username" binding:"required,max=32"`       // 登录用户名 1-32个字符
	Password string `json:"password" form:"password" binding:"required,min=6,max=32"` // 登录密码 6-32个字符
}

type UserInfoReq struct {
	User_ID string `json:"user_id" form:"user_id" binding:"required"` // 用户id
	Token   string `json:"token" form:"token" binding:"required"`     // 用户鉴权token
}
