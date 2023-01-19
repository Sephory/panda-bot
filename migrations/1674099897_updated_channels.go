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

		// add
		new_is_active := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "pwbclfhm",
			"name": "is_active",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_active)
		collection.Schema.AddField(new_is_active)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("pwbclfhm")

		return dao.SaveCollection(collection)
	})
}
