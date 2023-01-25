package twitch

import "github.com/sephory/panda-bot/pkg/chat"

var _ chat.ChatChannel = &TwitchChatChannel{}

type TwitchChatChannel struct {
	name            string
	chatConnection  *twitchChatConnection
	eventConnection *twitchEventConnection
	events          chan interface{}
}

func NewTwitchChatChannel(name string,
	chatConnection *twitchChatConnection,
	eventConnection *twitchEventConnection,

) *TwitchChatChannel {
	return &TwitchChatChannel{
		name:       name,
		chatConnection: chatConnection,
		events:     make(chan interface{}),
	}
}

func (c *TwitchChatChannel) GetName() string {
	return c.name
}

func (c *TwitchChatChannel) GetEvents() chan interface{} {
	go c.getChat()
	return c.events
}

func (c *TwitchChatChannel) SendMessage(message string, options interface{}) {
	c.chatConnection.message(c.name, message)
}

func (c *TwitchChatChannel) getChat() {
	channelMessages := c.chatConnection.joinChannel(c.name)
	if channelMessages == nil {
		return
	}
	for message := range channelMessages {
		event := message.toChatEvent()
		if event != nil {
			c.events <- event
		}
	}
}
