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
			"id": "7dfuhxpskzfvofw",
			"created": "2023-01-22 03:35:07.395Z",
			"updated": "2023-01-22 03:35:07.395Z",
			"name": "poll_votes",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "wgw2t5e2",
					"name": "poll_option_id",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "jamvcqjnxjp9yl9",
						"cascadeDelete": true
					}
				},
				{
					"system": false,
					"id": "wxotxds1",
					"name": "channel_id",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "vijqso7vo52sspt",
						"cascadeDelete": false
					}
				},
				{
					"system": false,
					"id": "hd72sxxy",
					"name": "voter",
					"type": "text",
					"required": false,
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

		collection, err := dao.FindCollectionByNameOrId("7dfuhxpskzfvofw")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
