package database

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
)

type Database struct {
	dao *daos.Dao
}

func New(pocketbase *pocketbase.PocketBase) *Database {
	database := &Database{}
	pocketbase.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		database.dao = e.App.Dao()
		return nil
	})
	return database
}
