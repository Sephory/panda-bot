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
			"id": "ytzns3z96dn6dnw",
			"created": "2023-01-21 21:32:43.947Z",
			"updated": "2023-01-21 21:32:43.947Z",
			"name": "channel_commands",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "lsyyhjg8",
					"name": "channel_id",
					"type": "relation",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"collectionId": "vijqso7vo52sspt",
						"cascadeDelete": true
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

		collection, err := dao.FindCollectionByNameOrId("ytzns3z96dn6dnw")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
