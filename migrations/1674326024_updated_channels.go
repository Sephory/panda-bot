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

		// update
		edit_service := &schema.SchemaField{}
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
		}`), edit_service)
		collection.Schema.AddField(edit_service)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		// update
		edit_service := &schema.SchemaField{}
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
					"Twitch"
				]
			}
		}`), edit_service)
		collection.Schema.AddField(edit_service)

		return dao.SaveCollection(collection)
	})
}
