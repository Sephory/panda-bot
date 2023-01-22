package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_POLLS = "polls"

var _ models.Model = &Poll{}

type Poll struct {
	models.BaseModel
	UserId   string
	Name     string
	Prompt   string
	IsOpen   bool
	IsActive bool
}

func (p *Poll) TableName() string {
	return TABLE_POLLS
}

func pollQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Poll{})
}

func (d *Database) GetActivePollsForUser(userId string) []*Poll {
	polls := []*Poll{}
	pollQuery(d.dao).
		Where(dbx.HashExp{"user_id": userId, "is_active": true}).
		All(&polls)
	return polls
}

func (d *Database) GetOldestActivePollForUser(userId string) *Poll {
	poll := &Poll{}
	pollQuery(d.dao).
		Where(dbx.HashExp{"user_id": userId, "is_active": true}).
		OrderBy("created").
		Limit(1).
		One(poll)
	return poll
}

func (d *Database) GetPollForUserByName(userId string, name string, activeOnly bool) *Poll {
	poll := &Poll{}
	query := pollQuery(d.dao).
		Where(dbx.HashExp{"user_id": userId, "name": name}).
		OrderBy("created DESC").
		Limit(1)
	if activeOnly {
		query = query.AndWhere(dbx.HashExp{"is_active": true})
	}
	query.One(poll)
	return poll
}

func (d *Database) GetMostRecentPollForUser(userId string) *Poll {
	poll := &Poll{}
	pollQuery(d.dao).
		Where(dbx.HashExp{"user_id": userId}).
		OrderBy("created DESC").
		Limit(1).
		One(poll)
	return poll
}

func (d *Database) SavePoll(poll *Poll) {
	d.dao.Save(poll)
}
