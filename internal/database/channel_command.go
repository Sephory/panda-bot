package database

import "github.com/pocketbase/pocketbase/models"

const TABLE_CHANNEL_COMMANDS = "channel_commands"

var _ models.Model = &ChannelCommand{}

type ChannelCommand struct {
	models.BaseModel
	ChannelId string
	CommandId string
}

func (cc *ChannelCommand) TableName() string {
	return TABLE_CHANNEL_COMMANDS
}

func (d *Database) SetDefaultChannelCommands(channelId string) {
	commands := []*Command{}
	commandQuery(d.dao).All(&commands)
	for _, cmd := range commands {
		channelCommand := &ChannelCommand{
			ChannelId: channelId,
			CommandId: cmd.Id,
		}
		err := d.dao.Save(channelCommand)
		if err != nil {
			break
		}
	}
}
