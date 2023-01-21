//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pocketbase/pocketbase"
	"github.com/sephory/panda-bot/internal/database"
	"github.com/sephory/panda-bot/internal/panda"
)

func InitializePanda() (*panda.Panda, error) {
	wire.Build(
		ReadConfig,
		GetChatClients,
		pocketbase.New,
		panda.New,
		panda.NewBot,
		database.New,
	)
	return &panda.Panda{}, nil
}
