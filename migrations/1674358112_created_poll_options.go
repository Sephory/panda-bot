package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "jamvcqjnxjp9yl9",
			"created": "2023-01-22 03:28:32.406Z",
			"updated": "2023-01-22 03:28:32.406Z",
			"name": "poll_options",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "3qrf5bca",
					"name": "poll_id",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "rp5cdp90kkpke2f",
						"cascadeDelete": true
					}
				},
				{
					"system": false,
					"id": "vmzr2pel",
					"name": "title",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("jamvcqjnxjp9yl9")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
