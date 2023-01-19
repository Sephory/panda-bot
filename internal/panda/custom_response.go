package panda

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/sephory/panda-bot/internal/database"
)

func getCustomResponse(dao *daos.Dao, channelName string, command Command) string {
	response := database.GetCustomResponseForChannel(dao, channelName, command.CommandText)

	if response.IsModOnly && !command.Event.User.IsMod {
		return ""
	}

	return response.Response
}
