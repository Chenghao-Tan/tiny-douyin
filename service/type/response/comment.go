package response

type CommentResp struct {
	Status
	Comment Comment `json:"comment"` // 评论成功返回评论内容，不需要重新拉取整个列表
}

type CommentListResp struct {
	Status
	Comment_List []Comment `json:"comment_list"` // 评论列表
}
