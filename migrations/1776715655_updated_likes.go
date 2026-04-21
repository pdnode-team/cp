package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2190274710")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_unique_user_like` + "`" + ` ON ` + "`" + `likes` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `target_id` + "`" + `,\n  ` + "`" + `target_collection` + "`" + `\n)",
				"CREATE INDEX ` + "`" + `idx_target_lookup` + "`" + ` ON ` + "`" + `likes` + "`" + ` (\n  ` + "`" + `target_id` + "`" + `,\n  ` + "`" + `target_collection` + "`" + `\n)"
			]
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2190274710")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_unique_user_like` + "`" + ` ON ` + "`" + `likes` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `target_id` + "`" + `,\n  ` + "`" + `target_collection` + "`" + `\n)"
			]
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
