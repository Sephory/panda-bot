package panda

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/sephory/panda-bot/internal/database"
	"github.com/sephory/panda-bot/pkg/chat"
	"github.com/sephory/panda-bot/pkg/chat/twitch"
)

const (
	CLIENT_TWITCH = "Twitch"
)

type pandaBot struct {
	dao *daos.Dao
	clients       map[string]chat.ChatClient
	commandPrefix byte
}

func newBot(config Config) *pandaBot {
	twitchClient := twitch.NewTwitchClient(config.Twitch)
	twitchClient.Connect()

	clients := map[string]chat.ChatClient{
		"Twitch": twitchClient,
	}
	return &pandaBot{
		clients:       clients,
		commandPrefix: config.CommandPrefix[0],
	}
}

func (bot *pandaBot) start(dao *daos.Dao) {
	bot.dao = dao
	bot.joinActive()
}

func (bot *pandaBot) joinActive() {
	channels := database.GetJoinedChannels(bot.dao)
	for _, channel := range channels {
		go bot.join(channel.Service, channel.Name)
	}
}

func (bot *pandaBot) join(serviceName string, channelName string) {
	client := bot.clients[serviceName]
	channel := client.JoinChannel(channelName)
	log.Printf("Joined %s channel %s", serviceName, channelName)
	for event := range channel.GetEvents() {
		log.Printf("(%s) %s: %s", channel.GetName(), event.User.DisplayName, event.Message)
		if event.Message[0] == bot.commandPrefix {
			command := NewCommand(event)
			bot.handleCommand(command, channel)
		}
	}
	log.Printf("Left %s channel %s", serviceName, channelName)
}

func (bot *pandaBot) leave(serviceName string, channelName string) {
	client := bot.clients[serviceName]
	client.LeaveChannel(channelName)
}

func (bot *pandaBot) onChannelSaved(channelId string) {
	channel := database.FindChannelById(bot.dao, channelId)
	if (channel.IsJoined) {
		go bot.join(channel.Service, channel.Name)
	} else {
		bot.leave(channel.Service, channel.Name)
	}
}

func (bot *pandaBot) onChannelDeleted(channelId string) {
	channel := database.FindChannelById(bot.dao, channelId)
	bot.leave(channel.Service, channel.Name)
}

func (bot *pandaBot) handleCommand(command Command, channel chat.ChatChannel) {
	var response string
	switch command.CommandType {
	case HelloWorld:
		response = fmt.Sprintf("Hello, %s!", command.Event.User.DisplayName)
	case RollDice:
		response = rollDice(command)
	default:
		response = getCustomResponse(bot.dao, channel.GetName(), command)
	}
	if (response != "") {
		log.Printf("SEND (%s): %s", channel.GetName(), response)
		channel.SendMessage(response)
	}
}

