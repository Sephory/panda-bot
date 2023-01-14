package panda

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
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
		go app.bot.start(e.App.Dao())
		return nil
	})

	return app.pocketbase.Start()
}
