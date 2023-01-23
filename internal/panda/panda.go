package panda

import (
	"github.com/pocketbase/pocketbase"
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

func (p *Panda) Start() error {
	migratecmd.MustRegister(p.pocketbase, p.pocketbase.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	p.pocketbase.OnBeforeServe().Add(p.onBeforeServe)

	p.pocketbase.OnModelAfterCreate().Add(p.onModelAfterCreate)

	p.pocketbase.OnModelAfterUpdate().Add(p.onModelAfterUpdate)

	p.pocketbase.OnModelBeforeDelete().Add(p.onModelBeforeDelete)

	return p.pocketbase.Start()
}
