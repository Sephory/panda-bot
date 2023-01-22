package panda

import (
	"fmt"
	"log"

	"github.com/sephory/panda-bot/internal/database"
	"github.com/sephory/panda-bot/pkg/chat"
)

type Bot struct {
	db      *database.Database
	clients map[string]chat.ChatClient
}

func NewBot(database *database.Database, clients ...chat.ChatClient) *Bot {
	clientMap := map[string]chat.ChatClient{}
	for _, client := range clients {
		clientMap[client.GetName()] = client
	}

	return &Bot{
		db:      database,
		clients: clientMap,
	}
}

func (bot *Bot) Start() {
	bot.joinActive()
}

func (bot *Bot) joinActive() {
	channels := bot.db.GetJoinedChannels()
	for _, channel := range channels {
		go bot.Join(channel)
	}
}

func (bot *Bot) Join(channel *database.Channel) {
	service := bot.db.FindServiceById(channel.ServiceId)
	client := bot.clients[service.Name]
	chatChannel := client.JoinChannel(channel.Name)
	log.Printf("Joined %s channel %s", service.Name, channel.Name)
	for e := range chatChannel.GetEvents() {
		switch event := e.(type) {
		case chat.Message:
			log.Printf("(%s) %s: %s", chatChannel.GetName(), event.User.DisplayName, event.Message)
			if event.Message[0] == '!' {
				command := NewCommand(event)
				bot.handleCommand(command, chatChannel, channel)
			}
		case chat.UserJoin:
			log.Printf("(%s) %s joined the channel", chatChannel.GetName(), event.User.Username)

		}
	}
	log.Printf("Left %s channel %s", service.Name, channel.Name)
}

func (bot *Bot) Leave(channel *database.Channel) {
	service := bot.db.FindServiceById(channel.ServiceId)
	client := bot.clients[service.Name]
	client.LeaveChannel(channel.Name)
}

func (bot *Bot) onChannelSaved(channelId string) {
	channel := bot.db.FindChannelById(channelId)
	if channel.IsJoined {
		go bot.Join(channel)
	} else {
		bot.Leave(channel)
	}
}

func (bot *Bot) onChannelDeleted(channelId string) {
	channel := bot.db.FindChannelById(channelId)
	bot.Leave(channel)
}

func (bot *Bot) handleCommand(command Command, chatChannel chat.ChatChannel, channelData *database.Channel) {
	var response string
	if !bot.db.IsCommandEnabledOnChannel(chatChannel.GetName(), command.CommandText) {
		response = bot.getCustomResponse(chatChannel.GetName(), command)
	} else {
		switch command.CommandType {
		case HelloWorld:
			response = fmt.Sprintf("Hello, %s!", command.Event.User.DisplayName)
		case RollDice:
			response = bot.rollDice(command)
		case Poll:
			response = bot.poll(command, channelData)
		case Vote:
			response = bot.vote(command, channelData)
		case Results:
			response = bot.results(command, channelData)
		}
	}
	if response != "" {
		log.Printf("SEND (%s): %s", chatChannel.GetName(), response)
		chatChannel.SendMessage(response)
	}
}
