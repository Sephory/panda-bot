package twitch

import "github.com/sephory/panda-bot/pkg/chat"

var _ chat.ChatChannel = &TwitchChatChannel{}

type TwitchChatChannel struct {
	name       string
	connection *twitchChatConnection
	events     chan interface{}
}

func NewTwitchChatChannel(name string, connection *twitchChatConnection) *TwitchChatChannel {
	return &TwitchChatChannel{
		name:       name,
		connection: connection,
		events:     make(chan interface{}),
	}
}

func (c *TwitchChatChannel) GetName() string {
	return c.name
}

func (c *TwitchChatChannel) GetEvents() chan interface{} {
	return c.events
}

func (c *TwitchChatChannel) SendMessage(message string, options interface{}) {
	c.connection.message(c.name, message)
}
