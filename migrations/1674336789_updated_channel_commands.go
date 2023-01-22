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

		collection, err := dao.FindCollectionByNameOrId("ytzns3z96dn6dnw")
		if err != nil {
			return err
		}

		// add
		new_command_id := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "8m9j6vyz",
			"name": "command_id",
			"type": "relation",
			"required": true,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"collectionId": "xhsbh0ot74znhin",
				"cascadeDelete": false
			}
		}`), new_command_id)
		collection.Schema.AddField(new_command_id)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("ytzns3z96dn6dnw")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("8m9j6vyz")

		return dao.SaveCollection(collection)
	})
}
