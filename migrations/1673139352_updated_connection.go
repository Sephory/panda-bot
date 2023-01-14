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
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		collection.Name = "channels"

		// update
		edit_user_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "3olwlrfw",
			"name": "user_id",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "_pb_users_auth_",
				"cascadeDelete": true
			}
		}`), edit_user_id)
		collection.Schema.AddField(edit_user_id)

		// update
		edit_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "yijhcqpq",
			"name": "name",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		collection.Name = "connection"

		// update
		edit_user_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_user_id)
		collection.Schema.AddField(edit_user_id)

		// update
		edit_name := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_name)
		collection.Schema.AddField(edit_name)

		return dao.SaveCollection(collection)
	})
}
