// main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cp-website/ent"
	"cp-website/handler" // 引入刚才写的 handler 包

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
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
	e.Use(middleware.Recover())

	// 使用 handler 包里的验证器
	e.Validator = &handler.CustomValidator{Validator: validator.New()}

	// 全局错误处理器
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := err.Error()

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
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello World"})
	})

	e.GET("/cp", handler.GetAllCP(client))
	e.POST("/cp", handler.CreateCP(client))
	e.GET("/cp/:id", handler.GetCP(client))
	e.PUT("/cp/:id", handler.UpdateCP(client))
	e.DELETE("/cp/:id", handler.DeleteCP(client))

	e.Logger.Fatal(e.Start(":8000"))
}
