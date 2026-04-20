package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2190274710")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "select3466706339",
			"maxSelect": 1,
			"name": "target_collection",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"cps",
				"characters"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_2190274710")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "select3466706339",
			"maxSelect": 1,
			"name": "target_collection",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"cps",
				"ccharacters"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
