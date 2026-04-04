package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// RateLimit 全局限流（按客户端 IP，见 echo.Context.RealIP）。
//
// 环境变量（可选）：
//   - RATE_LIMIT_RPS：平均每秒请求数，默认 30
//   - RATE_LIMIT_BURST：突发桶容量，默认 60
//
// 部署在反向代理后请设置环境变量 TRUSTED_PROXY_CIDRS（见 ConfigureIPExtractor），否则请保持未设置以使用直连 IP。
func RateLimit() echo.MiddlewareFunc {
	rps := 30.0
	burst := 60
	if s := os.Getenv("RATE_LIMIT_RPS"); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil && v > 0 {
			rps = v
		}
	}
	if s := os.Getenv("RATE_LIMIT_BURST"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			burst = v
		}
	}

	store := emiddleware.NewRateLimiterMemoryStoreWithConfig(emiddleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(rps),
		Burst:     burst,
		ExpiresIn: 5 * time.Minute,
	})

	return emiddleware.RateLimiterWithConfig(emiddleware.RateLimiterConfig{
		Store: store,
		Skipper: func(c echo.Context) bool {
			if c.Request().Method == http.MethodOptions {
				return true
			}
			switch c.Request().URL.Path {
			case "/", "/health":
				return true
			default:
				return false
			}
		},
		DenyHandler: func(c echo.Context, _ string, _ error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests, "Too many requests")
		},
	})
}
