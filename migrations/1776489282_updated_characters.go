package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3298390430")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_0KNEl6uZEK` + "`" + ` ON ` + "`" + `characters` + "`" + ` (` + "`" + `name` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_pPvOZMUzIt` + "`" + ` ON ` + "`" + `characters` + "`" + ` (` + "`" + `owner` + "`" + `)"
			]
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation1407448347")

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3298390430")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_0KNEl6uZEK` + "`" + ` ON ` + "`" + `characters` + "`" + ` (` + "`" + `name` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_pPvOZMUzIt` + "`" + ` ON ` + "`" + `characters` + "`" + ` (` + "`" + `owner` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_VLFJGGryYx` + "`" + ` ON ` + "`" + `characters` + "`" + ` (` + "`" + `cps` + "`" + `)"
			]
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
			"cascadeDelete": true,
			"collectionId": "pbc_1556475754",
			"hidden": false,
			"id": "relation1407448347",
			"maxSelect": 999,
			"minSelect": 0,
			"name": "cps",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
