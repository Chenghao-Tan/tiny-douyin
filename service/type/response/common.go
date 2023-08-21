package response

// 状态信息
type Status struct {
	Status_Code int    `json:"status_code"` // 状态码，0-成功，其他值-失败
	Status_Msg  string `json:"status_msg"`  // 返回状态描述
}

// 用户信息/视频作者信息/评论用户信息
type User struct {
	ID               uint   `json:"id"`               // 用户id
	Name             string `json:"name"`             // 用户名称
	Follow_Count     uint   `json:"follow_count"`     // 关注总数
	Follower_Count   uint   `json:"follower_count"`   // 粉丝总数
	Is_Follow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Avatar           string `json:"avatar"`           // 用户头像
	Background_Image string `json:"background_image"` // 用户个人页顶部大图
	Signature        string `json:"signature"`        // 个人简介
	Total_Favorited  string `json:"total_favorited"`  // 获赞数量
	Work_Count       uint   `json:"work_count"`       // 作品数
	Favorite_Count   uint   `json:"favorite_count"`   // 喜欢数
}

// 视频信息
type Video struct {
	ID             uint   `json:"id"`             // 视频唯一标识
	Author         User   `json:"author"`         // 视频作者信息
	Play_URL       string `json:"play_url"`       // 视频播放地址
	Cover_URL      string `json:"cover_url"`      // 视频封面地址
	Favorite_Count uint   `json:"favorite_count"` // 视频的点赞总数
	Comment_Count  uint   `json:"comment_count"`  // 视频的评论总数
	Is_Favorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title          string `json:"title"`          // 视频标题
}

// 评论信息
type Comment struct {
	ID          uint   `json:"id"`          // 评论id
	User        User   `json:"user"`        // 评论用户信息
	Content     string `json:"content"`     // 评论内容
	Create_Date string `json:"create_date"` // 评论发布日期，格式 mm-dd
}

// 聊天信息
type Message struct {
	ID           uint   `json:"id"`           // 消息id
	To_User_ID   uint   `json:"to_user_id"`   // 消息接收者id
	From_User_ID uint   `json:"from_user_id"` // 消息发送者id
	Content      string `json:"content"`      // 消息内容
	Create_Time  uint   `json:"create_time"`  // 消息发送时间 API文档有误 实为毫秒时间戳
}

// 用户(好友)信息
type FriendUser struct {
	User
	Message  string `json:"message"` // 和该好友的最新聊天消息
	Msg_Type uint   `json:"msgType"` // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}
