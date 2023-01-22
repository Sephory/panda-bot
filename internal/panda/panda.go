package panda

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/sephory/panda-bot/internal/database"
	_ "github.com/sephory/panda-bot/migrations"
)

type Panda struct {
	bot        *Bot
	pocketbase *pocketbase.PocketBase
	database   *database.Database
}

func New(bot *Bot, pocketbase *pocketbase.PocketBase, database *database.Database) *Panda {
	return &Panda{
		bot:        bot,
		pocketbase: pocketbase,
		database:   database,
	}

}

func (app *Panda) Start() error {
	migratecmd.MustRegister(app.pocketbase, app.pocketbase.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	app.pocketbase.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		app.bot.Start()
		return nil
	})

	app.pocketbase.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == database.TABLE_CHANNELS {
			channelId := e.Model.GetId()
			app.database.SetDefaultChannelCommands(channelId)
			app.bot.onChannelSaved(channelId)
		}
		return nil
	})

	app.pocketbase.OnModelAfterUpdate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == database.TABLE_CHANNELS {
			app.bot.onChannelSaved(e.Model.GetId())
		}
		return nil
	})

	app.pocketbase.OnModelBeforeDelete().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == database.TABLE_CHANNELS {
			app.bot.onChannelDeleted(e.Model.GetId())
		}
		return nil
	})

	return app.pocketbase.Start()
}
