package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_CHANNELS = "channels"

var _ models.Model = &Channel{}

type Channel struct {
	models.BaseModel
	UserId    string
	ServiceId string
	Name      string
	IsJoined  bool
}

func (c *Channel) TableName() string {
	return TABLE_CHANNELS
}

func channelQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Channel{})
}

func (d *Database) GetJoinedChannels() []*Channel {
	channels := []*Channel{}
	channelQuery(d.dao).
		Where(dbx.HashExp{"is_joined": true}).
		All(&channels)
	return channels
}

func (d *Database) FindChannelById(id string) *Channel {
	channel := &Channel{}
	channelQuery(d.dao).
		Where(dbx.HashExp{"id": id}).
		Limit(1).
		One(channel)
	return channel
}
