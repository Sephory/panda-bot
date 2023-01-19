package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qnxrfd2ae89pbdg")
		if err != nil {
			return err
		}

		collection.Name = "custom_responses"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qnxrfd2ae89pbdg")
		if err != nil {
			return err
		}

		collection.Name = "custom_response"

		return dao.SaveCollection(collection)
	})
}
