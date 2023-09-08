package request

type CommentReq struct {
	Token        string `json:"token" form:"token" binding:"required,jwt"`                                                    // 用户鉴权token
	Video_ID     uint   `json:"video_id" form:"video_id" binding:"required,min=1"`                                            // 视频id
	Action_Type  int    `json:"action_type" form:"action_type" binding:"required,min=1,max=2"`                                // 1-发布评论，2-删除评论
	Comment_Text string `json:"comment_text" form:"comment_text" binding:"required_if=Action_Type 1,omitempty,min=1,max=256"` // 可选参数，用户填写的评论内容，在action_type=1的时候使用
	Comment_ID   uint   `json:"comment_id" form:"comment_id" binding:"required_if=Action_Type 2,omitempty,min=1"`             // 可选参数，要删除的评论id，在action_type=2的时候使用
}

type CommentListReq struct {
	Token    string `json:"token" form:"token" binding:"omitempty,jwt"`        // 用户鉴权token API文档有误 应为可选参数
	Video_ID uint   `json:"video_id" form:"video_id" binding:"required,min=1"` // 视频id
}
