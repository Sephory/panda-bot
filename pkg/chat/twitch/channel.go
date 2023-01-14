package twitch

import "github.com/sephory/panda-bot/pkg/chat"

type TwitchChatChannel struct {
	name       string
	connection *twitchChatConnection
	events     chan chat.ChatEvent
}

func NewTwitchChatChannel(name string, connection *twitchChatConnection) *TwitchChatChannel {
	return &TwitchChatChannel{
		name:       name,
		connection: connection,
		events:     make(chan chat.ChatEvent),
	}
}

func (c *TwitchChatChannel) GetName() string {
	return c.name
}

func (c *TwitchChatChannel) GetEvents() chan chat.ChatEvent {
	return c.events
}

func (c *TwitchChatChannel) SendMessage(message string) {
	c.connection.sendMessage(c.name, message)
}
