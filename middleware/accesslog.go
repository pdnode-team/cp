package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// SlogAccessLog 使用结构化日志记录请求（含 X-Request-ID），路径不含 Query 以免泄露敏感参数。
func SlogAccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			slog.Info("http_request",
				"request_id", rid,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
			return err
		}
	}
}
