// handler/cp.go
package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"cp-website/ent"
	"cp-website/ent/cp"
	"cp-website/ent/tag"
	"cp-website/ent/user"
	"cp-website/model"

	"github.com/labstack/echo/v4"
)

// CheckOwnership 检查是否有权操作该资源
func CheckOwnership(c echo.Context, ownerID int64) error {
	currentUser := c.Get("user").(*ent.User)
	isAdmin := c.Get("isAdmin").(bool)

	if currentUser.ID != ownerID && !isAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "No right to operate other people's resources")
	}
	return nil
}

// 此函数负责查重、创建新标签，并返回所有相关标签的 ID 列表
func syncTags(ctx context.Context, client *ent.Client, names []string, dbUser *ent.User) ([]int64, error) {
	// 去重并剔除空字符串
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

	// 全局查找数据库中已有的标签 (不再限制 Owner)
	existingTags, err := client.Tag.Query().
		Where(tag.NameIn(nameList...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]int64)
	var tagIDs []int64
	for _, t := range existingTags {
		existingMap[t.Name] = t.ID
		tagIDs = append(tagIDs, t.ID)
	}

	// 找出需要新创建的标签
	var newCreates []*ent.TagCreate
	for _, n := range nameList {
		if _, exists := existingMap[n]; !exists {
			// 虽然是公共标签，但我们可以记录一下是谁第一个创建的
			newCreates = append(newCreates, client.Tag.Create().SetName(n).SetOwner(dbUser))
		}
	}

	// 4. 批量创建
	if len(newCreates) > 0 {
		createdTags, err := client.Tag.CreateBulk(newCreates...).Save(ctx)
		if err != nil {
			return nil, err
		}
		for _, t := range createdTags {
			tagIDs = append(tagIDs, t.ID)
		}
	}

	slog.Info("syncTags", "tagIDs", tagIDs)

	return tagIDs, nil
}

// GetAllCP 处理 GET /cp
func GetAllCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		cps, err := client.CP.Query().WithTags().All(ctx)
		if err != nil {
			return err
		}

		// 定义一个专为列表展示优化的响应结构
		type CPListItem struct {
			*ent.CP
			LikeCount int `json:"like_count"`
		}

		var result []CPListItem
		for _, item := range cps {
			// 在列表中逐个统计点赞数 (如果数据量极大，Ent 有更高级的 sql.Annotation 聚合语法，但目前这样写足够清晰)
			count, _ := item.QueryLikedByUsers().Count(ctx)
			result = append(result, CPListItem{
				CP:        item,
				LikeCount: count,
			})
		}

		return Success(c, result, 200)
	}
}

// GetCP 处理 GET /cp/:id
func GetCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := c.Param("id")
		val, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}

		// 1. 获取 CP 基本信息和它的标签
		dbCP, err := client.CP.Query().
			Where(cp.ID(val)).
			WithTags().
			Only(ctx)
		if err != nil {
			return err
		}

		// 2. 🌟 高性能获取总点赞数
		// QueryLikedByUsers() 会生成一条专门去中间表 COUNT 的 SQL 语句，不拉取具体用户数据
		likeCount, err := dbCP.QueryLikedByUsers().Count(ctx)
		if err != nil {
			return err
		}

		// 3. 🌟 判断当前看这个接口的用户，有没有点赞过这个 CP
		var isLiked bool

		// 尝试从 Context 中获取当前登录的用户（因为这个接口可能允许未登录查看，所以要做安全判断）
		if userVal := c.Get("user"); userVal != nil {
			dbUser := userVal.(*ent.User)
			// 查询当前 CP 的点赞列表中，是否存在当前用户的 ID
			isLiked, _ = dbCP.QueryLikedByUsers().
				Where(user.ID(dbUser.ID)).
				Exist(ctx)
		}

		// 4. 将 CP 数据、点赞数、当前用户的点赞状态 重新组装成 JSON 返回
		return Success(c, map[string]interface{}{
			"cp":         dbCP,      // 包含原本的 ID、名称、标签等
			"like_count": likeCount, // 总点赞数
			"is_liked":   isLiked,   // 当前用户是否已赞 (true/false)
		}, 200)
	}
}

// CreateCP 处理 POST /cp
func CreateCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		req := new(model.CPReq) // 注意这里引用了 model 包
		dbUser := c.Get("user").(*ent.User)

		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		exists, err := client.CP.Query().Where(cp.Name(req.Name)).Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			return echo.NewHTTPError(http.StatusBadRequest, "The CP name already exists.")
		}

		tagIDs, err := syncTags(ctx, client, req.TagNames, dbUser)
		if err != nil {
			return err
		}

		builder := client.CP.Create().
			SetName(req.Name).
			SetCategory(req.Category).
			SetOwner(dbUser).
			AddTagIDs(tagIDs...)

		if req.Link != nil {
			builder.SetLink(*req.Link)
		}

		newCP, err := builder.Save(ctx)
		if err != nil {
			return err
		}

		result, _ := client.CP.Query().Where(cp.ID(newCP.ID)).WithTags().Only(ctx)
		slog.Info("CreateCP", "result", result)
		return Success(c, result, 201)
	}
}

// DeleteCP 处理 DELETE /cp/:id
func DeleteCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := ParseID(c.Param("id"))

		dbCP, err := client.CP.Query().Where(cp.ID(id)).WithOwner().Only(ctx)

		if err != nil {
			return err // 让全局处理器返回 404 或 500，防止宕机
		}

		if err := CheckOwnership(c, dbCP.Edges.Owner.ID); err != nil {
			return err // 直接返回 403
		}

		if err := client.CP.DeleteOneID(dbCP.ID).Exec(ctx); err != nil {
			return err
		}
		slog.Info("DeleteCP", "dbCP", dbCP)
		return Success(c, dbCP, 200)
	}
}

// UpdateCP 处理 PUT /cp/:id
func UpdateCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := ParseID(c.Param("id"))
		dbUser := c.Get("user").(*ent.User)

		req := new(model.CPReq)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		dbCP, err := client.CP.Query().Where(cp.ID(id)).WithOwner().Only(ctx)

		if err != nil {
			return err // 让全局处理器返回 404 或 500，防止宕机
		}

		if err := CheckOwnership(c, dbCP.Edges.Owner.ID); err != nil {
			return err // 直接返回 403
		}

		tagIDs, err := syncTags(ctx, client, req.TagNames, dbUser)
		if err != nil {
			return err
		}

		builder := client.CP.UpdateOneID(id).
			SetName(req.Name).
			SetCategory(req.Category).
			ClearTags().
			AddTagIDs(tagIDs...)

		if req.Link != nil {
			builder.SetLink(*req.Link)
		} else {
			builder.ClearLink()
		}

		updatedCP, err := builder.Save(ctx)
		if err != nil {
			return err
		}
		slog.Info("UpdateCP", "updatedCP", updatedCP)
		result, _ := client.CP.Query().Where(cp.ID(updatedCP.ID)).WithTags().Only(ctx)
		return Success(c, result, 200)
	}
}
