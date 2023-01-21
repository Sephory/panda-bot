package panda

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/sephory/panda-bot/internal/database"
	_ "github.com/sephory/panda-bot/migrations"
	"github.com/sephory/panda-bot/pkg/chat/twitch"
)

type Config struct {
	CommandPrefix string `yaml:"command_prefix"`
	Twitch        twitch.TwitchClientConfiguration
}

type App struct {
	bot        *pandaBot
	pocketbase *pocketbase.PocketBase
}

func New(config Config) *App {
	return &App{
		bot:        newBot(config),
		pocketbase: pocketbase.New(),
	}

}

func (app *App) Start() error {
	migratecmd.MustRegister(app.pocketbase, app.pocketbase.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	app.pocketbase.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		app.bot.start(e.App.Dao())
		return nil
	})

	app.pocketbase.OnModelAfterCreate().Add(func(e *core.ModelEvent) error {
		if e.Model.TableName() == database.TABLE_CHANNELS {
			app.bot.onChannelSaved(e.Model.GetId())
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
