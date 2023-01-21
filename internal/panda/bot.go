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
	dao           *daos.Dao
	clients       map[string]chat.ChatClient
	commandPrefix byte
}

func newBot(config Config) *pandaBot {
	twitchClient := twitch.NewTwitchClient(config.Twitch)
	err := twitchClient.Connect()
	if err != nil {
		panic(err)
	}

	clients := map[string]chat.ChatClient{
		CLIENT_TWITCH: twitchClient,
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
	for e := range channel.GetEvents() {
		switch event := e.(type) {
		case chat.Message:
			log.Printf("(%s) %s: %s", channel.GetName(), event.User.DisplayName, event.Message)
			if event.Message[0] == bot.commandPrefix {
				command := NewCommand(event)
				bot.handleCommand(command, channel)
			}
		case chat.UserJoin:
			log.Printf("(%s) %s joined the channel", channel.GetName(), event.User.Username)

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
	if channel.IsJoined {
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
	if response != "" {
		log.Printf("SEND (%s): %s", channel.GetName(), response)
		channel.SendMessage(response)
	}
}
