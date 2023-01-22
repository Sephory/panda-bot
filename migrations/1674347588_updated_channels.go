package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("fhx18zbq")

		// add
		new_service_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "3eekjfez",
			"name": "service_id",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "3lnpxdbiatdnebg",
				"cascadeDelete": false
			}
		}`), new_service_id)
		collection.Schema.AddField(new_service_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		// add
		del_service := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "fhx18zbq",
			"name": "service",
			"type": "select",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"Twitch",
					"YouTube"
				]
			}
		}`), del_service)
		collection.Schema.AddField(del_service)

		// remove
		collection.Schema.RemoveField("3eekjfez")

		return dao.SaveCollection(collection)
	})
}
