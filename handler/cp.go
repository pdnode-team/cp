// handler/cp.go
package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"cp-website/ent"
	"cp-website/ent/cp"
	"cp-website/ent/tag"
	"cp-website/model"

	"github.com/labstack/echo/v4"
)

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

// GetAllCP 处理 GET /cp
func GetAllCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		cps, err := client.CP.Query().WithTags().All(c.Request().Context())
		if err != nil {
			return err
		}
		return Success(c, cps, 200)
	}
}

// GetCP 处理 GET /cp/:id
func GetCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		val, err := strconv.ParseInt(id, 10, 64) // 直接用 ParseInt
		if err != nil {
			return err
		}

		cpResult, err := client.CP.Query().Where(cp.ID(val)).WithTags().First(c.Request().Context())
		if err != nil {
			return err
		}
		return Success(c, cpResult, 200)
	}
}

// CreateCP 处理 POST /cp
func CreateCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		req := new(model.CPReq) // 注意这里引用了 model 包

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

		tagIDs, err := syncTags(ctx, client, req.TagNames)
		if err != nil {
			return err
		}

		builder := client.CP.Create().
			SetName(req.Name).
			SetCategory(req.Category).
			AddTagIDs(tagIDs...)

		if req.Link != nil {
			builder.SetLink(*req.Link)
		}

		newCP, err := builder.Save(ctx)
		if err != nil {
			return err
		}

		result, _ := client.CP.Query().Where(cp.ID(newCP.ID)).WithTags().Only(ctx)
		return Success(c, result, 201)
	}
}

// DeleteCP 处理 DELETE /cp/:id
func DeleteCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := ParseID(c.Param("id"))

		dbCP, err := client.CP.Query().Where(cp.IDEQ(id)).WithTags().Only(ctx)
		if err != nil {
			return err
		}

		if err := client.CP.DeleteOneID(dbCP.ID).Exec(ctx); err != nil {
			return err
		}
		return Success(c, dbCP, 200)
	}
}

// UpdateCP 处理 PUT /cp/:id
func UpdateCP(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		id := ParseID(c.Param("id"))

		req := new(model.CPReq)
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

		result, _ := client.CP.Query().Where(cp.ID(updatedCP.ID)).WithTags().Only(ctx)
		return Success(c, result, 200)
	}
}
