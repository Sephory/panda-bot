package twitch

import (

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
		channels: make(map[string]*TwitchChatChannel),
	}
	return &client
}

// JoinChannel implements chat.ChatClient
func (c *TwitchClient) JoinChannel(channelName string) chat.ChatChannel {
	c.api.subscribe(channelName, c.eventConnection.getSessionId())

	channel := NewTwitchChatChannel(channelName, c.chatConnection, c.eventConnection)
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
