// handler/common.go
package handler

import (
	"cp-website/util"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator 包装 validator 以适配 Echo
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// Success 统一的成功响应格式 (注意首字母大写，才能被其他包调用)
func Success(c echo.Context, data interface{}, status_code int) error {
	if status_code == 0 {
		status_code = http.StatusOK
	}
	return c.JSON(status_code, map[string]interface{}{
		"status": "ok",
		"data":   data,
	})
}

// ParseID 辅助函数：安全转换路由参数 ID (首字母大写)
func ParseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, util.ErrInvalidID
	}
	return id, nil
}
