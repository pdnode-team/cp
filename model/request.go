// model/request.go
package model

// CPReq 用于接收前端 POST/PUT 数据的结构体
type CPReq struct {
	Name     string   `json:"name" validate:"required,min=1,max=50"`
	Category string   `json:"category" validate:"required,min=2,max=20"`
	Link     *string  `json:"link" validate:"omitempty,url"`
	TagNames []string `json:"tag_names" validate:"required,min=1,max=15"`
}

// CommentReq 用于接收前端发表评论的请求
type CommentReq struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
	// 可选参数：如果是回复别人，前端就传这个 parent_id；如果是直接评论 CP，就不传
	ParentID *int64 `json:"parent_id" validate:"omitempty,min=1"` // 雪花 ID，勿用 max 限制数值上限
}
