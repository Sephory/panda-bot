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
			"id": "3lnpxdbiatdnebg",
			"created": "2023-01-22 00:31:46.930Z",
			"updated": "2023-01-22 00:31:46.930Z",
			"name": "services",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "l6oeiyjv",
					"name": "name",
					"type": "text",
					"required": true,
					"unique": true,
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

		collection, err := dao.FindCollectionByNameOrId("3lnpxdbiatdnebg")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
