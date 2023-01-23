package panda

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/sephory/panda-bot/internal/database"
)

func (p *Panda) onModelAfterCreate(e *core.ModelEvent) error {
	if e.Model.TableName() == database.TABLE_CHANNELS {
		channelId := e.Model.GetId()
		p.database.SetDefaultChannelCommands(channelId)
		p.bot.onChannelSaved(channelId)
	}
	return nil
}

func (p *Panda) onModelAfterUpdate(e *core.ModelEvent) error {
	if e.Model.TableName() == database.TABLE_CHANNELS {
		p.bot.onChannelSaved(e.Model.GetId())
	}
	return nil
}

func (p *Panda) onModelBeforeDelete(e *core.ModelEvent) error {
	if e.Model.TableName() == database.TABLE_CHANNELS {
		p.bot.onChannelDeleted(e.Model.GetId())
	}
	return nil
}
