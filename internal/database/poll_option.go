package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_POLL_OPTIONS = "poll_options"

var _ models.Model = &PollOption{}

type PollOption struct {
	models.BaseModel
	PollId    string
	Text      string
	ChannelId string
	CreatedBy string
}

func (po *PollOption) TableName() string {
	return TABLE_POLL_OPTIONS
}

func pollOptionQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&PollOption{})
}

func (d *Database) GetPollOptionsForPoll(pollId string) []*PollOption {
	pollOptions := []*PollOption{}
	pollOptionQuery(d.dao).
		Where(dbx.HashExp{"poll_id": pollId}).
		All(&pollOptions)
	return pollOptions
}

func (d *Database) SavePollOption(pollOption *PollOption) {
	d.dao.Save(pollOption)
}
