package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := `{
			"createRule": "@request.auth.id != \"\" && owner = @request.auth.id",
			"deleteRule": "@request.auth.id != \"\" && owner = @request.auth.id",
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text1579384326",
					"max": 0,
					"min": 0,
					"name": "name",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"convertURLs": false,
					"hidden": false,
					"id": "editor1843675174",
					"maxSize": 0,
					"name": "description",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "editor"
				},
				{
					"hidden": false,
					"id": "file3309110367",
					"maxSelect": 3,
					"maxSize": 0,
					"mimeTypes": [
						"image/png",
						"image/jpeg",
						"image/webp"
					],
					"name": "image",
					"presentable": false,
					"protected": false,
					"required": true,
					"system": false,
					"thumbs": [],
					"type": "file"
				},
				{
					"cascadeDelete": false,
					"collectionId": "_pb_users_auth_",
					"hidden": false,
					"id": "relation3479234172",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "owner",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "relation"
				},
				{
					"hidden": false,
					"id": "json3504296006",
					"maxSize": 0,
					"name": "tag_names",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3332085495",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_1556475754",
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_XG4ctO5tml` + "`" + ` ON ` + "`" + `cps` + "`" + ` (` + "`" + `name` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_Wwu23uTke6` + "`" + ` ON ` + "`" + `cps` + "`" + ` (` + "`" + `owner` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_0Fjn0gP3R4` + "`" + ` ON ` + "`" + `cps` + "`" + ` (\n  ` + "`" + `name` + "`" + `,\n  ` + "`" + `owner` + "`" + `\n)"
			],
			"listRule": "",
			"name": "cps",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id != \"\" && owner = @request.auth.id && @request.body.owner:changed = false",
			"viewRule": ""
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1556475754")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
