package twitch

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/sephory/panda-bot/pkg/chat"
)

const twitch_chat_url = "wss://irc-ws.chat.twitch.tv"
const twitch_chat_auth_timeout = 500

type twitchChatConnection struct {
	connection     *chat.WebsocketConnection
	authenticating sync.Mutex
	authenticated  bool
	token          string
	username       string
	messages       map[string]chan TwitchMessage
	log            log.Interface
}

func newTwitchChatConnection(token string, username string) *twitchChatConnection {
	chatUrl, _ := url.Parse(twitch_chat_url)
	return &twitchChatConnection{
		connection: chat.NewWebsocketConnection(chatUrl),
		token:      token,
		username:   username,
		messages:   make(map[string]chan TwitchMessage),
		log:        log.WithField("module", "twitch_chat_connection"),
	}
}

func (c *twitchChatConnection) awaitAuthentication() error {
	c.authenticating.Lock()
	defer c.authenticating.Unlock()
	if c.authenticated {
		return nil
	}
	authSuccess := c.authenticate()
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_chat_auth_timeout)
		timeout <- errors.New("Failed to authenticate with Twitch chat")
	}()
	select {
	case success := <-authSuccess:
		c.authenticated = success
		go c.listen()
		return nil
	case err := <-timeout:
		c.authenticated = false
		return err
	}
}

func (c *twitchChatConnection) disconnect() {
	c.connection.SendMessage(chat.WebsocketMessage{
		MessageType: websocket.CloseMessage,
		Bytes:       websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	})
}

func (c *twitchChatConnection) authenticate() chan bool {
	result := make(chan bool, 1)
	go func() {
		ircMessages, err := c.connection.GetMessages()
		if err != nil {
			result <- false
			return
		}
		for {
			messages := c.parseTwitchMessages(<-ircMessages)
			for _, m := range messages {
				if m.MessageType == AuthSuccess {
					result <- true
					return
				}
			}
		}
	}()
	c.send("CAP", "REQ", ":twitch.tv/membership twitch.tv/tags twitch.tv/commands")
	c.send("PASS", fmt.Sprintf("oauth:%s", c.token))
	c.send("NICK", c.username)
	return result
}

func (c *twitchChatConnection) joinChannel(channelName string) chan TwitchMessage {
	err := c.awaitAuthentication()
	if err != nil {
		c.log.Error(err.Error())
		return nil
	}
	c.send("JOIN", "#"+channelName)
	channelEvents := make(chan TwitchMessage)
	c.messages[channelName] = channelEvents
	return channelEvents
}

func (c *twitchChatConnection) leaveChannel(channelName string) {
	close(c.messages[channelName])
	delete(c.messages, channelName)
	c.send("PART", "#"+channelName)
	if len(c.messages) == 0 {
		// TODO: disconnect
	}
}

func (c *twitchChatConnection) message(channel, message string) {
	c.send("PRIVMSG", "#"+channel, ":"+message)
}

func (c *twitchChatConnection) send(command string, params ...string) {
	message := []byte(command + " ")
	message = append(message, strings.Join(params, " ")...)
	c.connection.SendMessage(chat.WebsocketMessage{
		MessageType: websocket.TextMessage,
		Bytes:       message,
	})
}

func (c *twitchChatConnection) listen() {
	messages, err := c.connection.GetMessages()
	if err != nil {
		return
	}
	for ircBatch := range messages {
		twitchMessages := c.parseTwitchMessages(ircBatch)
		for _, m := range twitchMessages {
			if m.MessageType == Ping {
				c.send("PONG", m.Params...)
			} else {
				channelName := m.getChannel()
				if channelMessages, ok := c.messages[channelName]; ok {
					channelMessages <- m
				}
			}
		}
	}
}

func (c *twitchChatConnection) parseTwitchMessages(ircBatch chat.WebsocketMessage) []TwitchMessage {
	messageStrings := strings.Split(string(ircBatch.Bytes), "\r\n")
	messages := []TwitchMessage{}
	for _, s := range messageStrings {
		if s == "" {
			continue
		}
		messages = append(messages, NewChatMessage(s))
	}
	return messages
}
