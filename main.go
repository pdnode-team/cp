package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/osutils"
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

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
