package main

import (
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/osutils"

	_ "pdnode.com/x/cp/migrations"
)

//go:embed all:pb_public
var embeddedFiles embed.FS

var Version = "untracked"

func main() {
	app := pocketbase.New()

	err := godotenv.Load()
	if err != nil {
		app.Logger().Warn("Warning: .env file not found, using system environment variables", "err", err)
	}

	app.RootCmd.Version = Version

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the IsProbablyGoRun check is to enable it only during development)
		Automigrate: osutils.IsProbablyGoRun(),
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		Setup(se.App)

		se.Router.POST("/api/{type}/{id}/toggle-like", func(e *core.RequestEvent) error {
			requestCollection := e.Request.PathValue("type")
			id := e.Request.PathValue("id")

			if requestCollection != "cps" && requestCollection != "characters" {
				return e.Error(400, "Invalid request data", map[string]validation.Error{
					"type": validation.NewError("validation_collection_type", "Invalid collection type."),
				})
			}

			collection, err := app.FindCollectionByNameOrId(requestCollection)
			if err != nil {
				return e.InternalServerError("Internal Server Error", map[string]any{"message": "Cannot find collection", "err": err})
			}

			record, err := app.FindRecordById(requestCollection, id)

			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return e.InternalServerError("Internal Server Error", map[string]any{
						"message": "Cannot find record",
						"err":     err,
					})
				}
			}

			if record != nil {

				err = app.Delete(record)
				if err != nil {
					return e.InternalServerError("Internal Server Error", map[string]any{
						"message": "Cannot delete record",
						"err":     err,
					})
				}

				return e.JSON(200, map[string]any{
					"message": "Success",
					"like":    false,
					"data":    record,
				})
			}

			record = core.NewRecord(collection)

			record.Set("target_id", id)
			record.Set("target_collection", requestCollection)
			record.Set("owner", e.Auth.Id)

			err = app.Save(record)

			if err != nil {
				return e.InternalServerError("Internal Server Error", map[string]any{
					"message": "Cannot save record",
					"err":     err,
				})
			}

			return e.JSON(200, map[string]any{
				"message": "Success",
				"like":    true,
				"data":    record,
			})

		}).Bind(apis.RequireAuth())

		publicFS, err := fs.Sub(embeddedFiles, "pb_public")
		if err != nil {
			return err
		}
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(publicFS, true))

		return se.Next()
	})

	app.OnRecordCreateRequest("cps").BindFunc(validateCharactersOwnership)
	app.OnRecordUpdateRequest("cps").BindFunc(validateCharactersOwnership)

	app.OnRecordCreateRequest().BindFunc(validateIdImmutable)

	app.OnRecordCreateRequest().BindFunc(restrictToAuth)
	app.OnRecordUpdateRequest().BindFunc(restrictToAuth)
	app.OnRecordDeleteRequest().BindFunc(restrictToAuth)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// 权限校验中间件
func restrictToAuth(e *core.RecordRequestEvent) error {
	if e.Auth != nil || e.Record.Collection().IsAuth() {
		return e.Next()
	}
	return apis.NewForbiddenError("Unauthorized", nil)
}

// 阻止自定义IDs
func validateIdImmutable(e *core.RecordRequestEvent) error {
	// Superuser 绕过
	if e.Auth != nil && e.Auth.IsSuperuser() {
		return e.Next()
	}

	// 检查用户有没有传入id
	if e.Record.Id != "" {
		return e.BadRequestError("Invalid request data", map[string]validation.Error{
			"id": validation.NewError("validation_id_immutable", "Custom record IDs are not allowed."),
		})
	}

	return e.Next()
}

// 校验 CP 中的Characters是否属于当前操作者
func validateCharactersOwnership(e *core.RecordRequestEvent) error {

	if e.Auth == nil {
		return e.UnauthorizedError("Unauthorized", nil)
	}

	// Superuser 绕过
	if e.Auth != nil && e.Auth.IsSuperuser() {
		return e.Next()
	}

	// 获取当前记录中关联的角色 IDs (假设字段名是 "members")
	characterIds := e.Record.GetStringSlice("characters")
	if len(characterIds) == 0 {
		return e.Next()
	}
	// 将 slice 转为 interface slice 用于 dbx 查询
	vals := make([]any, len(characterIds))
	for i, v := range characterIds {
		vals[i] = v
	}

	// 查询这些角色中，是否存在不属于当前用户的记录
	var count int
	err := e.App.DB().
		Select("count(*)").
		From("characters").
		Where(dbx.And(
			dbx.In("id", vals...),
			dbx.Not(dbx.HashExp{"owner": e.Auth.Id}), // 找出 owner 不是当前用户的
		)).
		Row(&count)

	if err != nil {
		return e.InternalServerError("Database query failed", err)
	}

	// 如果发现有角色不属于该用户，拦截请求
	if count > 0 {
		return e.Error(403, "Illegal association", map[string]validation.Error{
			"characters": validation.NewError("invalid_character_owner", "One or more selected characters do not belong to you."),
		})
	}

	return e.Next()
}
