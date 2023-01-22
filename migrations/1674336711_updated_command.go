package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("xhsbh0ot74znhin")
		if err != nil {
			return err
		}

		collection.Name = "commands"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("xhsbh0ot74znhin")
		if err != nil {
			return err
		}

		collection.Name = "command"

		return dao.SaveCollection(collection)
	})
}
