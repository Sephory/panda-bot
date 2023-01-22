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

		collection, err := dao.FindCollectionByNameOrId("jamvcqjnxjp9yl9")
		if err != nil {
			return err
		}

		// add
		new_channel_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "1t09tup7",
			"name": "channel_id",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "vijqso7vo52sspt",
				"cascadeDelete": false
			}
		}`), new_channel_id)
		collection.Schema.AddField(new_channel_id)

		// add
		new_created_by := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "quw3dfkw",
			"name": "created_by",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_created_by)
		collection.Schema.AddField(new_created_by)

		// update
		edit_text := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "vmzr2pel",
			"name": "text",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_text)
		collection.Schema.AddField(edit_text)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("jamvcqjnxjp9yl9")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("1t09tup7")

		// remove
		collection.Schema.RemoveField("quw3dfkw")

		// update
		edit_text := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_text)
		collection.Schema.AddField(edit_text)

		return dao.SaveCollection(collection)
	})
}
