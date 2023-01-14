package database

import "github.com/pocketbase/pocketbase/daos"

type Database struct {
	dao *daos.Dao
}

func New(dao *daos.Dao) *Database {
	return &Database{
		dao: dao,
	}
}

func (db *Database) GetAllChannels() []Channel {
	var channels []Channel
	db.dao.DB().NewQuery("SELECT * FROM channels").All(&channels)
	return channels
}


