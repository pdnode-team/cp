// main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cp-website/auth"
	"cp-website/ent"
	"cp-website/handler" // 引入刚才写的 handler 包
	"cp-website/middleware"
	"cp-website/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// listenAddr resolves the HTTP bind address from LISTEN_ADDR or PORT.
// LISTEN_ADDR takes precedence (e.g. ":8000", "0.0.0.0:8000", "[::]:8000").
// If only PORT is set (e.g. "8080"), the server listens on 0.0.0.0:PORT.
// Default: 0.0.0.0:8000.
func listenAddr() string {
	if a := strings.TrimSpace(os.Getenv("LISTEN_ADDR")); a != "" {
		return a
	}
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		if strings.HasPrefix(p, ":") {
			return "0.0.0.0" + p
		}
		return "0.0.0.0:" + p
	}
	return "0.0.0.0:8000"
}

func main() {

	// 启动时只读取一次 Secret
	fileBytes, err := os.ReadFile("secret.json")
	if err != nil {
		log.Fatalf("无法读取 secret.json: %v", err)
	}
	var zitadelConfig auth.ZitadelConfig
	if err := json.Unmarshal(fileBytes, &zitadelConfig); err != nil {
		log.Fatalf("解析 secret.json 失败: %v", err)
	}

	client, err := ent.Open("sqlite3", "file:database.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }
	// fmt.Println("数据库表创建成功！")

	e := echo.New()
	middleware.ConfigureIPExtractor(e)

	bodyLimit := os.Getenv("HTTP_BODY_LIMIT")
	if bodyLimit == "" {
		bodyLimit = "512K"
	}
	e.Use(emiddleware.BodyLimit(bodyLimit))

	e.Use(emiddleware.RequestID())
	e.Use(middleware.SlogAccessLog())
	e.Use(emiddleware.Secure())
	e.Use(emiddleware.Recover())
	e.Use(middleware.RateLimit())
	e.Use(emiddleware.CORSWithConfig(emiddleware.CORSConfig{
		AllowOrigins:     []string{"https://cp.p67.click", "http://localhost:8000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXRequestID},
		ExposeHeaders:    []string{echo.HeaderXRequestID},
		AllowCredentials: true, // 若前端要带 Cookie / Authorization，常需 true，且 AllowOrigins 不能是 "*"
	}))

	// 使用 handler 包里的验证器
	e.Validator = &handler.CustomValidator{Validator: validator.New()}

	// 全局错误处理器（5xx 对外统一文案，细节写入 slog）
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		rid := c.Response().Header().Get(echo.HeaderXRequestID)
		code := http.StatusInternalServerError
		msg := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if code >= http.StatusInternalServerError {
				slog.Error("http error",
					"request_id", rid,
					"code", code,
					"message", he.Message,
					"internal", he.Internal,
				)
				msg = "Internal Server Error"
			} else {
				msg = fmt.Sprintf("%v", he.Message)
			}
		} else if ent.IsNotFound(err) {
			code = http.StatusNotFound
			msg = "Data not found"
		} else if errors.Is(err, strconv.ErrSyntax) {
			code = http.StatusBadRequest
			msg = "Please provide the correct number."
		} else if ent.IsConstraintError(err) {
			code = http.StatusBadRequest
			msg = "Unique constraint failed (e.g., name already exists)"
		} else if errors.Is(err, util.ErrInvalidID) {
			code = http.StatusBadRequest
			msg = "Invalid ID"
		} else if err != nil {
			slog.Error("unhandled error", "request_id", rid, "error", err)
		}

		_ = c.JSON(code, map[string]interface{}{
			"status": "error",
			"msg":    msg,
		})
	}

	// 注册路由，直接调用 handler 包里的方法
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "CP Website API is running!"})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
	})

	api := e.Group("/cp")
	// 给这个组挂载你的鉴权中间件
	api.Use(middleware.AuthMiddleware(client, zitadelConfig))

	// 下面所有的路由都会自动先经过 AuthMiddleware 校验！
	api.GET("", handler.GetAllCP(client))
	api.POST("", handler.CreateCP(client))
	api.GET("/:id", handler.GetCP(client))
	api.PUT("/:id", handler.UpdateCP(client))
	api.DELETE("/:id", handler.DeleteCP(client))

	api.POST("/:id/like", handler.ToggleLike(client))
	api.POST("/:id/comment", handler.CreateComment(client))
	api.GET("/:id/comments", handler.GetComments(client))

	addr := listenAddr()
	slog.Info("server listening", "addr", addr)
	e.Logger.Fatal(e.Start(addr))
}
