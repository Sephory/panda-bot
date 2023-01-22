package panda

func (bot *Bot) getCustomResponse(channelName string, command Command) string {
	response := bot.db.GetCustomResponseForChannel(channelName, command.CommandText)

	if response.IsModOnly && !command.Event.User.IsMod {
		return ""
	}

	return response.Response
}
