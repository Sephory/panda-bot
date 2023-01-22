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

		collection, err := dao.FindCollectionByNameOrId("xhsbh0ot74znhin")
		if err != nil {
			return err
		}

		// update
		edit_title := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "l9oyqzqe",
			"name": "title",
			"type": "text",
			"required": true,
			"unique": true,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_title)
		collection.Schema.AddField(edit_title)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("xhsbh0ot74znhin")
		if err != nil {
			return err
		}

		// update
		edit_title := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "l9oyqzqe",
			"name": "name",
			"type": "text",
			"required": true,
			"unique": true,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_title)
		collection.Schema.AddField(edit_title)

		return dao.SaveCollection(collection)
	})
}
