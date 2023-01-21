package panda

import (
	"fmt"
	"log"

	"github.com/sephory/panda-bot/internal/database"
	"github.com/sephory/panda-bot/pkg/chat"
)

type Bot struct {
	database *database.Database
	clients  map[string]chat.ChatClient
}

func NewBot(database *database.Database, clients ...chat.ChatClient) *Bot {
	clientMap := map[string]chat.ChatClient{}
	for _, client := range clients {
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		clientMap[client.GetName()] = client
	}

	return &Bot{
		database: database,
		clients:  clientMap,
	}
}

func (bot *Bot) Start() {
	bot.joinActive()
}

func (bot *Bot) joinActive() {
	channels := bot.database.GetJoinedChannels()
	for _, channel := range channels {
		go bot.Join(channel.Service, channel.Name)
	}
}

func (bot *Bot) Join(serviceName string, channelName string) {
	client := bot.clients[serviceName]
	channel := client.JoinChannel(channelName)
	log.Printf("Joined %s channel %s", serviceName, channelName)
	for e := range channel.GetEvents() {
		switch event := e.(type) {
		case chat.Message:
			log.Printf("(%s) %s: %s", channel.GetName(), event.User.DisplayName, event.Message)
			if event.Message[0] == '!' {
				command := NewCommand(event)
				bot.handleCommand(command, channel)
			}
		case chat.UserJoin:
			log.Printf("(%s) %s joined the channel", channel.GetName(), event.User.Username)

		}
	}
	log.Printf("Left %s channel %s", serviceName, channelName)
}

func (bot *Bot) Leave(serviceName string, channelName string) {
	client := bot.clients[serviceName]
	client.LeaveChannel(channelName)
}

func (bot *Bot) onChannelSaved(channelId string) {
	channel := bot.database.FindChannelById(channelId)
	if channel.IsJoined {
		go bot.Join(channel.Service, channel.Name)
	} else {
		bot.Leave(channel.Service, channel.Name)
	}
}

func (bot *Bot) onChannelDeleted(channelId string) {
	channel := bot.database.FindChannelById(channelId)
	bot.Leave(channel.Service, channel.Name)
}

func (bot *Bot) handleCommand(command Command, channel chat.ChatChannel) {
	var response string
	switch command.CommandType {
	case HelloWorld:
		response = fmt.Sprintf("Hello, %s!", command.Event.User.DisplayName)
	case RollDice:
		response = bot.rollDice(command)
	default:
		response = bot.getCustomResponse(channel.GetName(), command)
	}
	if response != "" {
		log.Printf("SEND (%s): %s", channel.GetName(), response)
		channel.SendMessage(response)
	}
}
