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
			"id": "vijqso7vo52sspt",
			"created": "2023-01-07 21:49:45.623Z",
			"updated": "2023-01-07 21:49:45.623Z",
			"name": "connection",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "3olwlrfw",
					"name": "user",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": true
					}
				},
				{
					"system": false,
					"id": "fhx18zbq",
					"name": "service",
					"type": "select",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"values": [
							"Twitch"
						]
					}
				},
				{
					"system": false,
					"id": "yijhcqpq",
					"name": "channel",
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
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
