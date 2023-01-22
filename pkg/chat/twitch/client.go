package twitch

import (
	"log"

	"github.com/sephory/panda-bot/pkg/chat"
)

type TwitchClientConfiguration struct {
	Token    string
	Username string
	ClientId string `yaml:"client_id"`
}

var _ chat.ChatClient = &TwitchClient{}

type TwitchClient struct {
	chatConnection  *twitchChatConnection
	eventConnection *twitchEventConnection
	api             *twitchApi
	channels        map[string]*TwitchChatChannel
}

func New(config *TwitchClientConfiguration) *TwitchClient {
	client := TwitchClient{
		chatConnection:  newTwitchChatConnection(config.Token, config.Username),
		eventConnection: newTwitchEventConnection(),
		api:             newTwitchApi(config.ClientId, config.Token),
		channels:        map[string]*TwitchChatChannel{},
	}
	go client.listen()
	return &client
}

// JoinChannel implements chat.ChatClient
func (c *TwitchClient) JoinChannel(channelName string) chat.ChatChannel {
	c.chatConnection.joinChannel(channelName)
	c.api.subscribe(channelName, c.eventConnection.getSessionId())

	channel := NewTwitchChatChannel(channelName, c.chatConnection)
	c.channels[channelName] = channel
	return channel
}

// LeaveChannel implements chat.ChatClient
func (c *TwitchClient) LeaveChannel(channelName string) {
	if c.channels[channelName] == nil {
		return
	}
	c.chatConnection.leaveChannel(channelName)
	close(c.channels[channelName].events)
	delete(c.channels, channelName)
	if len(c.channels) == 0 {
		c.chatConnection.disconnect()
		c.eventConnection.disconnect()
	}
}

// GetName implements chat.ChatClient
func (c *TwitchClient) GetName() string {
	return "Twitch"
}

func (c *TwitchClient) listen() {
	for {
		select {
		case m := <-c.chatConnection.messages:
			event, err := m.toChatEvent()
			if err != nil {
				continue
			}
			if channel, ok := c.channels[getChannel(m)]; ok {
				channel.events <- event
			}
		case e := <-c.eventConnection.events:
			log.Println(e)
		}
	}
}

func getChannel(message TwitchMessage) string {
	switch message.MessageType {
	case ChatMessage:
		return message.Params[0][1:]
	case Join:
		return message.Params[0][1:]
	}
	return ""
}
