package middleware

import (
	"github.com/labstack/echo/v4"
)

// WriteAPIError 与全局 HTTPErrorHandler 一致：{"status":"error","msg":...}
func WriteAPIError(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]interface{}{
		"status": "error",
		"msg":    msg,
	})
}
