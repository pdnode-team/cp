package main

import (
	"embed"
	"io/fs"
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the IsProbablyGoRun check is to enable it only during development)
		Automigrate: osutils.IsProbablyGoRun(),
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		publicFS, err := fs.Sub(embeddedFiles, "pb_public")
		if err != nil {
			return err
		}
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(publicFS, true))

		return se.Next()
	})

	app.OnRecordCreateRequest("characters").BindFunc(validateCPsOwnership)
	app.OnRecordUpdateRequest("characters").BindFunc(validateCPsOwnership)

	app.OnRecordCreateRequest().BindFunc(validateIdImmutable)
	app.OnRecordUpdateRequest().BindFunc(validateIdImmutable)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// 阻止自定义IDs
func validateIdImmutable(e *core.RecordRequestEvent) error {
	// Superuser 绕过
	if e.Auth != nil && e.Auth.IsSuperuser() {
		return e.Next()
	}

	// 检查用户有没有传入id
	if e.Record.Id != "" {
		return apis.NewBadRequestError("Invalid request data", map[string]validation.Error{
			"id": validation.NewError("validation_id_immutable", "Custom record IDs are not allowed."),
		})
	}

	return e.Next()
}

// 通用校验逻辑
func validateCPsOwnership(e *core.RecordRequestEvent) error {
	// Superuser 绕过
	if e.Auth != nil && e.Auth.IsSuperuser() {
		return e.Next()
	}

	cpIds := e.Record.GetStringSlice("cps")
	if len(cpIds) == 0 {
		return e.Next()
	}

	vals := make([]any, len(cpIds))
	for i, v := range cpIds {
		vals[i] = v
	}

	var count int
	err := e.App.DB().
		Select("count(*)").
		From("cps").
		Where(dbx.And(
			dbx.In("id", vals...),
			dbx.Not(dbx.HashExp{"owner": e.Auth.Id}),
		)).
		Row(&count)

	if err != nil {
		return apis.NewInternalServerError("Database query failed", err)
	}

	if count > 0 {
		return apis.NewBadRequestError("Illegal association", map[string]validation.Error{
			"cps": validation.NewError("invalid_owner", "Includes a partner that doesn't belong to you."),
		})
	}

	return e.Next()
}
