package request

type CommentReq struct {
	Token        string `json:"token" form:"token" binding:"required"`             // 用户鉴权token
	Video_ID     string `json:"video_id" form:"video_id" binding:"required"`       // 视频id
	Action_Type  string `json:"action_type" form:"action_type" binding:"required"` // 1-发布评论，2-删除评论
	Comment_Text string `json:"comment_text,omitempty" form:"comment_text"`        // 可选参数，用户填写的评论内容，在action_type=1的时候使用
	Comment_ID   string `json:"comment_id,omitempty" form:"comment_id"`            // 可选参数，要删除的评论id，在action_type=2的时候使用
}

type CommentListReq struct {
	Token    string `json:"token,omitempty" form:"token"`                // 用户鉴权token API文档有误 应为可选参数
	Video_ID string `json:"video_id" form:"video_id" binding:"required"` // 视频id
}
