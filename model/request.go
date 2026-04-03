// model/request.go
package model

// CPReq 用于接收前端 POST/PUT 数据的结构体
type CPReq struct {
	Name     string   `json:"name" validate:"required,min=1"`
	Category string   `json:"category" validate:"required,min=2,max=20"`
	Link     *string  `json:"link"`
	TagNames []string `json:"tag_names" validate:"required,min=1"`
}
