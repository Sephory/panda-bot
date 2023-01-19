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

		// update
		edit_is_joined := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "pwbclfhm",
			"name": "is_joined",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), edit_is_joined)
		collection.Schema.AddField(edit_is_joined)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vijqso7vo52sspt")
		if err != nil {
			return err
		}

		// update
		edit_is_joined := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "pwbclfhm",
			"name": "is_active",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), edit_is_joined)
		collection.Schema.AddField(edit_is_joined)

		return dao.SaveCollection(collection)
	})
}
