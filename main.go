package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cp-website/ent"
	"cp-website/ent/cp"
	"cp-website/ent/tag"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// ---------------- 1. 定义 DTO 与校验器 ----------------
// 用于接收前端 POST/PUT 数据的结构体 (相当于你的 CPCreate)
type CPReq struct {
	Name     string   `json:"name" validate:"required,min=1"`
	Category string   `json:"category" validate:"required,min=2,max=20"`
	Link     *string  `json:"link"`
	TagNames []string `json:"tag_names" validate:"required,min=1"`
}

// CustomValidator 包装 validator 以适配 Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ---------------- 2. 辅助函数：处理标签逻辑 ----------------
// 此函数负责查重、创建新标签，并返回所有相关标签的 ID 列表
func syncTags(ctx context.Context, client *ent.Client, names []string) ([]int64, error) {
	// 1. 去重并剔除空字符串
	uniqueNames := make(map[string]bool)
	for _, n := range names {
		trimmed := strings.TrimSpace(n)
		if trimmed != "" {
			uniqueNames[trimmed] = true
		}
	}
	if len(uniqueNames) == 0 {
		return nil, echo.NewHTTPError(400, "At least one non-empty tag is required")
	}

	var nameList []string
	for n := range uniqueNames {
		nameList = append(nameList, n)
	}

	// 2. 查找数据库中已有的标签
	existingTags, err := client.Tag.Query().Where(tag.NameIn(nameList...)).All(ctx)
	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]int64) // 改为 int64
	var tagIDs []int64                    // 改为 int64
	for _, t := range existingTags {
		existingMap[t.Name] = t.ID
		tagIDs = append(tagIDs, t.ID) // 收集已有标签 ID
	}

	// 3. 找出需要新创建的标签
	var newCreates []*ent.TagCreate
	for _, n := range nameList {
		if _, exists := existingMap[n]; !exists {
			newCreates = append(newCreates, client.Tag.Create().SetName(n))
		}
	}

	// 4. 批量创建新标签并收集它们的 ID
	if len(newCreates) > 0 {
		createdTags, err := client.Tag.CreateBulk(newCreates...).Save(ctx)
		if err != nil {
			return nil, err
		}
		for _, t := range createdTags {
			tagIDs = append(tagIDs, t.ID)
		}
	}

	return tagIDs, nil
}

// ---------------- 3. 主程序 ----------------
func main() {
	// 初始化 Ent 客户端 (SQLite)
	client, err := ent.Open("sqlite3", "file:database.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// 自动运行迁移（建表）
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	fmt.Println("数据库表创建成功！")

	// 初始化 Echo
	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	// 全局错误处理器 (替代你的 universal_exception_handler)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		msg := err.Error()

		// 如果是 Echo 的 HTTPError，或者是实体不存在的错误
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

	// ---------------- 4. 路由定义 ----------------

	// 响应结构统一包装
	success := func(c echo.Context, data interface{}) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"data":   data,
		})
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello World"})
	})

	// GET /cp
	e.GET("/cp", func(c echo.Context) error {
		// 显式连表查询：WithTags()。查出来的结果就是纯数据，永远没有 Detached 错误！
		cps, err := client.CP.Query().WithTags().All(c.Request().Context())
		if err != nil {
			return err // 直接抛出，被全局错误处理器接管
		}
		return success(c, cps)
	})

	// POST /cp
	e.POST("/cp", func(c echo.Context) error {
		ctx := c.Request().Context()
		req := new(CPReq)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// 检查重名 (Ent 也能利用 ConstraintError 捕获，但显式检查更像你之前的逻辑)
		exists, err := client.CP.Query().Where(cp.Name(req.Name)).Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			return echo.NewHTTPError(http.StatusBadRequest, "The CP name already exists.")
		}

		// 处理标签并获取关联 ID
		tagIDs, err := syncTags(ctx, client, req.TagNames)
		if err != nil {
			return err
		}

		// 创建 CP 并关联标签
		builder := client.CP.Create().
			SetName(req.Name).
			SetCategory(req.Category).
			AddTagIDs(tagIDs...) // 魔法在此：直接通过 ID 关联！

		if req.Link != nil {
			builder.SetLink(*req.Link)
		}

		newCP, err := builder.Save(ctx)
		if err != nil {
			return err
		}

		// 显式查询一次完整的包含 Tag 的实体用于返回
		result, _ := client.CP.Query().Where(cp.ID(newCP.ID)).WithTags().Only(ctx)
		return success(c, result)
	})

	// DELETE /cp/:id
	e.DELETE("/cp/:id", func(c echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")

		// 1. 先查出完整数据（因为你的需求是要返回被删除的数据）
		// 在 Ent 中，查出来的数据就是纯内存 Struct，删了数据库的它也不会报错
		dbCP, err := client.CP.Query().Where(cp.IDEQ(parseID(id))).WithTags().Only(ctx)
		if err != nil {
			return err
		}

		// 2. 真实执行删除
		if err := client.CP.DeleteOneID(dbCP.ID).Exec(ctx); err != nil {
			return err
		}

		// 3. 返回刚才缓存的 struct 即可
		return success(c, dbCP)
	})

	// PUT /cp/:id
	e.PUT("/cp/:id", func(c echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")

		req := new(CPReq)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		tagIDs, err := syncTags(ctx, client, req.TagNames)
		if err != nil {
			return err
		}

		builder := client.CP.UpdateOneID(parseID(id)).
			SetName(req.Name).
			SetCategory(req.Category).
			ClearTags().         // 先清空旧的关联
			AddTagIDs(tagIDs...) // 添加新的关联

		if req.Link != nil {
			builder.SetLink(*req.Link)
		} else {
			builder.ClearLink()
		}

		updatedCP, err := builder.Save(ctx)
		if err != nil {
			return err
		}

		result, _ := client.CP.Query().Where(cp.ID(updatedCP.ID)).WithTags().Only(ctx)
		return success(c, result)
	})

	// GET /cp/:id
	e.GET("/cp/:id", func(c echo.Context) error {
		id := c.Param("id")
		val, err := strconv.Atoi(id)
		if err != nil {
			// 转换失败（比如里面有字母）
			return err
		}

		cp, err := client.CP.Query().Where(cp.ID(int64(val))).WithTags().First(c.Request().Context())
		if err != nil {
			return err // 直接抛出，被全局错误处理器接管
		}
		return success(c, cp)

	})

	e.Logger.Fatal(e.Start(":8000"))
}

// 辅助函数：安全转换路由参数 ID
func parseID(idStr string) int64 { // 返回类型改为 int64
	id, err := strconv.ParseInt(idStr, 10, 64) // 显式转换为 int64
	if err != nil {
		return 0 // 或者根据你的需要返回错误
	}
	return id
}
