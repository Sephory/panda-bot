package panda

func (bot *Bot) getCustomResponse(channelName string, command Command) string {
	response := bot.database.GetCustomResponseForChannel(channelName, command.CommandText)

	if response.IsModOnly && !command.Event.User.IsMod {
		return ""
	}

	return response.Response
}
