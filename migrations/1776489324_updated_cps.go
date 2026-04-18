package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1556475754")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"cascadeDelete": true,
			"collectionId": "pbc_3298390430",
			"hidden": false,
			"id": "relation975782158",
			"maxSelect": 2,
			"minSelect": 0,
			"name": "characters",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1556475754")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation975782158")

		return app.Save(collection)
	})
}
