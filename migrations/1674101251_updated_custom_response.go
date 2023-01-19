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

		collection, err := dao.FindCollectionByNameOrId("qnxrfd2ae89pbdg")
		if err != nil {
			return err
		}

		// add
		new_is_mod_only := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "d8limxm0",
			"name": "is_mod_only",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_mod_only)
		collection.Schema.AddField(new_is_mod_only)

		// add
		new_is_owner_only := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "k2q6gbmv",
			"name": "is_owner_only",
			"type": "bool",
			"required": false,
			"unique": false,
			"options": {}
		}`), new_is_owner_only)
		collection.Schema.AddField(new_is_owner_only)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qnxrfd2ae89pbdg")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("d8limxm0")

		// remove
		collection.Schema.RemoveField("k2q6gbmv")

		return dao.SaveCollection(collection)
	})
}
