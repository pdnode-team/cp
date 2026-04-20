package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3298390430")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("url3740358174")

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text3740358174",
			"max": 0,
			"min": 0,
			"name": "origin",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": true,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3298390430")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"exceptDomains": [],
			"hidden": false,
			"id": "url3740358174",
			"name": "origin",
			"onlyDomains": [],
			"presentable": false,
			"required": true,
			"system": false,
			"type": "url"
		}`)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text3740358174")

		return app.Save(collection)
	})
}
