// main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cp-website/auth"
	"cp-website/ent"
	"cp-website/handler" // 引入刚才写的 handler 包
	"cp-website/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

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

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	fmt.Println("数据库表创建成功！")

	e := echo.New()
	e.Use(emiddleware.Recover())

	// 使用 handler 包里的验证器
	e.Validator = &handler.CustomValidator{Validator: validator.New()}

	// 全局错误处理器
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			msg = fmt.Sprintf("%v", he.Message)
		} else if ent.IsNotFound(err) {
			code = http.StatusNotFound
			msg = "Data not found"
		} else if errors.Is(err, strconv.ErrSyntax) {
			code = http.StatusBadRequest
			msg = "Please provide the correct number."
		} else if ent.IsConstraintError(err) {
			code = http.StatusBadRequest
			msg = "Unique constraint failed (e.g., name already exists)"
		}

		c.JSON(code, map[string]interface{}{
			"status": "error",
			"msg":    msg,
		})
	}

	// 注册路由，直接调用 handler 包里的方法
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
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

	e.Logger.Fatal(e.Start(":8000"))
}
