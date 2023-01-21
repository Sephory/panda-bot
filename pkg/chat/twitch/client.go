package twitch

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sephory/panda-bot/pkg/chat"
)

const twitch_chat_auth_timeout = 500
const twitch_event_welcome_timeout = 500

type TwitchClientConfiguration struct {
	Token    string
	Username string
	ClientId string `yaml:"client_id"`
}

var _ chat.ChatClient = &TwitchClient{}

type TwitchClient struct {
	config          TwitchClientConfiguration
	chatConnection  *twitchChatConnection
	eventConnection *twitchEventConnection
	api             *twitchApi
	channels        map[string]*TwitchChatChannel
}

func NewTwitchClient(config TwitchClientConfiguration) *TwitchClient {
	client := TwitchClient{
		config:          config,
		chatConnection:  &twitchChatConnection{},
		eventConnection: &twitchEventConnection{},
		api:             &twitchApi{httpClient: http.DefaultClient, config: &config},
		channels:        map[string]*TwitchChatChannel{},
	}
	return &client
}

func (c *TwitchClient) Connect() error {
	if err := c.chatConnection.connect(); err != nil {
		return err
	}
	if err := c.awaitChatAuthentication(); err != nil {
		return err
	}
	if err := c.eventConnection.connect(); err != nil {
		return err
	}
	if err := c.awaitEventSession(); err != nil {
		return err
	}
	go c.listen()
	return nil
}

func (c *TwitchClient) Disconnect() {
	for _, channel := range c.channels {
		c.LeaveChannel(channel.name)
	}
	c.chatConnection.disconnect()
}

func (c *TwitchClient) JoinChannel(channelName string) chat.ChatChannel {
	c.chatConnection.joinChannel(channelName)
	c.api.subscribe(channelName, c.eventConnection.sessionId)
	channel := NewTwitchChatChannel(channelName, c.chatConnection)
	c.channels[channelName] = channel
	return channel
}

func (c *TwitchClient) LeaveChannel(channelName string) {
	if c.channels[channelName] == nil {
		return
	}
	c.chatConnection.leaveChannel(channelName)
	close(c.channels[channelName].events)
	delete(c.channels, channelName)
}

func (c *TwitchClient) awaitChatAuthentication() error {
	c.chatConnection.authenticate(c.config.Token, c.config.Username)
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_chat_auth_timeout)
		timeout <- errors.New("Failed to authenticate with Twitch chat")
	}()
	for {
		select {
		case m := <-c.chatConnection.messages:
			if m.MessageType == AuthSuccess {
				return nil
			}
		case err := <-timeout:
			return err
		}
	}
}

func (c *TwitchClient) awaitEventSession() error {
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_chat_auth_timeout)
		timeout <- errors.New("Twitch EventSub welcom timed out")
	}()
	for {
		select {
		case e := <-c.eventConnection.events:
			if welcome, ok := e.(*sessionWelcome); ok {
				c.eventConnection.sessionId = welcome.Session.Id
				return nil
			}
		case err := <-timeout:
			return err
		}
	}

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

func (c *TwitchClient) log(message ...interface{}) {
	log.Println("twitch_chat_client", message)
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
