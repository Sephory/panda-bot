package twitch

import (
	"errors"
	"log"
	"time"

	"github.com/sephory/panda-bot/pkg/chat"
)

const twitch_chat_auth_timeout = 500

type TwitchClientConfiguration struct {
	Token    string
	Username string
}

type TwitchClient struct {
	config     TwitchClientConfiguration
	connection *twitchChatConnection
	channels   map[string]*TwitchChatChannel
}

func NewTwitchClient(config TwitchClientConfiguration) *TwitchClient {
	client := TwitchClient{
		config:     config,
		connection: &twitchChatConnection{},
		channels:   map[string]*TwitchChatChannel{},
	}
	return &client
}

func (c *TwitchClient) Connect() error {
	err := c.connection.connect()
	if err == nil {
		err = c.authenticate()
	}
	if err == nil {
		go c.listen()
	}
	return err
}

func (c *TwitchClient) Disconnect() {
	for _, channel := range c.channels {
		c.LeaveChannel(channel.name)
	}
	c.connection.disconnect()
}

func (c *TwitchClient) JoinChannel(channelName string) chat.ChatChannel {
	c.connection.joinChannel(channelName)
	channel := NewTwitchChatChannel(channelName, c.connection)
	c.channels[channelName] = channel
	return channel
}

func (c *TwitchClient) LeaveChannel(channelName string) {
	if c.channels[channelName] == nil {
		return
	}
	c.connection.leaveChannel(channelName)
	close(c.channels[channelName].events)
	delete(c.channels, channelName)
}

func (c *TwitchClient) authenticate() error {
	c.connection.authenticate(c.config.Token, c.config.Username)
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_chat_auth_timeout)
		timeout <- errors.New("Failed to authenticate with Twitch chat")
	}()
	for {
		select {
		case m := <-c.connection.messages:
			if m.MessageType == AuthSuccess {
				return nil
			}
		case err := <-timeout:
			return err
		}
	}
}

func (c *TwitchClient) listen() {
	for {
		message := <-c.connection.messages
		switch message.MessageType {
		case ChatMessage:
			event := getChatEvent(message)
			if channel, ok := c.channels[getChannel(message)]; ok {
				channel.events <- event
			}
		}
	}
}

func getChatEvent(message TwitchMessage) chat.ChatEvent {
	return chat.ChatEvent{
		User:    getUserInfo(message),
		Message: message.Params[1],
	}
}

func (c *TwitchClient) log(message ...interface{}) {
	log.Println("twitch_chat_client", message)
}

func getUserInfo(message TwitchMessage) chat.UserInfo {
	userInfo := chat.UserInfo{}
	if message.MessageType == ChatMessage {
		userInfo.Username = message.Source.Name
		userInfo.DisplayName = message.Tags["display-name"][0]
		userInfo.IsMod = message.Tags["mod"][0] == "1"
		userInfo.IsSubscriber = message.Tags["subscriber"][0] == "1"
	}
	return userInfo
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
