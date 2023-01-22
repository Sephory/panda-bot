package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_COMMANDS = "commands"

var _ models.Model = &Command{}

type Command struct {
	models.BaseModel
	Title string
	Text  string
}

func (c *Command) TableName() string {
	return TABLE_COMMANDS
}

func commandQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Command{})
}

func (d *Database) IsCommandEnabledOnChannel(channelName string, commandText string) bool {
	rows, err := commandQuery(d.dao).
		InnerJoin(TABLE_CHANNEL_COMMANDS, dbx.NewExp("channel_commands.command_id = commands.id")).
		InnerJoin(TABLE_CHANNELS, dbx.NewExp("channels.id = channel_commands.channel_id")).
		Where(dbx.HashExp{
			"channels.name": channelName,
			"commands.text": commandText,
		}).
		Limit(1).
		Rows()
	if err != nil {
		return false
	}
	return rows.Next()
}
