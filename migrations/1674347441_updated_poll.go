package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("rp5cdp90kkpke2f")
		if err != nil {
			return err
		}

		collection.Name = "polls"

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("rp5cdp90kkpke2f")
		if err != nil {
			return err
		}

		collection.Name = "poll"

		return dao.SaveCollection(collection)
	})
}
