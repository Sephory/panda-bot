package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_CHANNELS = "channels"

type Channel struct {
	models.BaseModel
	UserId   string
	Service  string
	Name     string
	IsJoined bool
}

func (c *Channel) TableName() string {
	return TABLE_CHANNELS
}

func channelQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Channel{})
}

func GetJoinedChannels(dao *daos.Dao) []*Channel {
	channels := []*Channel{}
	channelQuery(dao).
		Where(dbx.HashExp{"is_joined": true}).
		All(&channels)
	return channels
}

func FindChannelById(dao *daos.Dao, id string) *Channel {
	channel := &Channel{}
	channelQuery(dao).
		Where(dbx.HashExp{"id": id}).
		Limit(1).
		One(channel)
	return channel
}
