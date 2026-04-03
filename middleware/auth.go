// middleware/auth.go
package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"cp-website/auth" // 确保这个路径对应你实际的 Zitadel 包
	"cp-website/ent"
	"cp-website/ent/user"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware 接收初始化好的 config，避免每次请求都读硬盘
func AuthMiddleware(client *ent.Client, config auth.ZitadelConfig) echo.MiddlewareFunc {
	apiURL := "https://auth.pdnode.com/oauth/v2/introspect"

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization Header"})
			}

			var token string
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Authorization Header Format"})
			}

			assertion, _ := auth.GenerateClientAssertion(config.ClientID, config.KeyID, config.Key, "https://auth.pdnode.com")
			data := url.Values{}
			data.Set("token", token)
			data.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
			data.Set("client_assertion", assertion)

			resp, err := http.PostForm(apiURL, data)
			if err != nil {
				slog.Error("请求 ZITADEL 失败", "error", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			}
			defer resp.Body.Close()

			var introspectionResult struct {
				Active   bool   `json:"active"`
				Sub      string `json:"sub"`
				Username string `json:"username"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&introspectionResult); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			}

			if !introspectionResult.Active {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token is invalid"})
			}

			// 🌟 核心逻辑：查询或创建用户实体
			ctx := c.Request().Context()
			dbUser, err := client.User.Query().Where(user.SubEQ(introspectionResult.Sub)).Only(ctx)
			if ent.IsNotFound(err) {
				dbUser, err = client.User.Create().SetSub(introspectionResult.Sub).Save(ctx)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
				}
			}

			// 🌟 魔法：把查出来的数据库 User 实体直接塞进 Context
			c.Set("user", dbUser)

			return next(c)
		}
	}
}
