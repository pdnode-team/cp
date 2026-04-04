// middleware/auth.go
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"cp-website/auth"
	"cp-website/ent"
	"cp-website/ent/user"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/singleflight"
)

func zitadelIntrospectURL() string {
	if u := strings.TrimSpace(os.Getenv("ZITADEL_INTROSPECT_URL")); u != "" {
		return u
	}
	return "https://auth.pdnode.com/oauth/v2/introspect"
}

func zitadelAudience() string {
	if a := strings.TrimSpace(os.Getenv("ZITADEL_AUDIENCE")); a != "" {
		return a
	}
	return "https://auth.pdnode.com"
}

func introspectCacheTTL() time.Duration {
	s := strings.TrimSpace(os.Getenv("ZITADEL_INTROSPECT_CACHE_TTL"))
	if s == "" {
		return 30 * time.Second
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0
	}
	return d
}

func zitadelHTTPClientTimeout() time.Duration {
	s := strings.TrimSpace(os.Getenv("ZITADEL_HTTP_TIMEOUT"))
	if s == "" {
		return 15 * time.Second
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 15 * time.Second
	}
	return d
}

func internalErr(err error) *echo.HTTPError {
	he := echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	he.Internal = err
	return he
}

func gatewayErr(err error) *echo.HTTPError {
	he := echo.NewHTTPError(http.StatusBadGateway, "Internal Server Error")
	he.Internal = err
	return he
}

func fetchIntrospection(ctx context.Context, httpClient *http.Client, apiURL string, form url.Values) (introspectionResult, error) {
	var zero introspectionResult
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return zero, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("introspect HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out introspectionResult
	if err := json.Unmarshal(body, &out); err != nil {
		return zero, fmt.Errorf("introspect decode: %w", err)
	}
	return out, nil
}

// AuthMiddleware 接收初始化好的 config；Zitadel 地址与 audience 可由环境变量覆盖，见 zitadel* 函数。
func AuthMiddleware(client *ent.Client, config auth.ZitadelConfig) echo.MiddlewareFunc {
	apiURL := zitadelIntrospectURL()
	audience := zitadelAudience()
	httpClient := &http.Client{Timeout: zitadelHTTPClientTimeout()}
	cache := newIntrospectCache(introspectCacheTTL())
	var sf singleflight.Group

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Response().Header().Get(echo.HeaderXRequestID)

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return WriteAPIError(c, http.StatusUnauthorized, "Missing authorization header")
			}

			var token string
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				return WriteAPIError(c, http.StatusBadRequest, "Invalid authorization header format")
			}

			assertion, err := auth.GenerateClientAssertion(config.ClientID, config.KeyID, config.Key, audience)
			if err != nil {
				slog.Error("zitadel client assertion", "error", err, "request_id", rid)
				return internalErr(err)
			}

			cacheKey := tokenCacheKey(token)
			if cache != nil {
				if res, ok := cache.get(cacheKey); ok {
					return finalizeAuth(c, client, res, next, rid)
				}
			}

			v, err, _ := sf.Do(cacheKey, func() (interface{}, error) {
				form := url.Values{}
				form.Set("token", token)
				form.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
				form.Set("client_assertion", assertion)
				res, err := fetchIntrospection(c.Request().Context(), httpClient, apiURL, form)
				if err != nil {
					return nil, err
				}
				if cache != nil {
					cache.set(cacheKey, res)
				}
				return res, nil
			})
			if err != nil {
				slog.Error("zitadel introspect", "error", err, "request_id", rid)
				return gatewayErr(err)
			}

			res := v.(introspectionResult)
			return finalizeAuth(c, client, res, next, rid)
		}
	}
}

func finalizeAuth(c echo.Context, client *ent.Client, introspectionResult introspectionResult, next echo.HandlerFunc, rid string) error {
	if !introspectionResult.Active {
		return WriteAPIError(c, http.StatusUnauthorized, "Invalid or expired token")
	}

	var roles []string
	isAdmin := false
	if introspectionResult.Scope != "" {
		if strings.Contains(introspectionResult.Scope, "role:admin") {
			isAdmin = true
			roles = append(roles, "admin")
		}
		if strings.Contains(introspectionResult.Scope, "role:moderator") {
			roles = append(roles, "moderator")
		}
	}

	ctx := c.Request().Context()
	dbUser, err := client.User.Query().Where(user.SubEQ(introspectionResult.Sub)).Only(ctx)
	if ent.IsNotFound(err) {
		dbUser, err = client.User.Create().SetSub(introspectionResult.Sub).Save(ctx)
		if err != nil {
			slog.Error("create user after introspect", "error", err, "request_id", rid)
			return internalErr(err)
		}
	} else if err != nil {
		slog.Error("load user", "error", err, "request_id", rid)
		return internalErr(err)
	}

	c.Set("user", dbUser)
	c.Set("roles", roles)
	c.Set("isAdmin", isAdmin)
	return next(c)
}
