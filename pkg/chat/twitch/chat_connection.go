package twitch

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const twitch_chat_host = "irc-ws.chat.twitch.tv:443"
const twitch_chat_auth_timeout = 500

type queueMessage struct {
	messageType int
	message     []byte
}

type twitchChatConnection struct {
	connection *websocket.Conn
	connecting bool
	token      string
	username   string
	messages   chan TwitchMessage
	sendQueue  chan queueMessage
}

func newTwitchChatConnection(token string, username string) *twitchChatConnection {
	return &twitchChatConnection{
		token:     token,
		username:  username,
		messages:  make(chan TwitchMessage),
		sendQueue: make(chan queueMessage),
	}
}

func (c *twitchChatConnection) connect() error {
	c.connecting = true
	defer func() { c.connecting = false }()

	url := url.URL{Scheme: "wss", Host: twitch_chat_host}
	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if connection != nil {
		c.connection = connection
		go c.writeMessages()
		err = c.awaitAuthentication(c.token, c.username)
		if err != nil {
			return err
		}
		go c.readMessages()
	}
	return err
}

func (c *twitchChatConnection) awaitAuthentication(token string, username string) error {
	c.authenticate(token, username)
	auth := make(chan error, 1)
	go func() {
		for {
			_, ircMessages, err := c.connection.ReadMessage()
			if err != nil {
				auth <- err
				return
			}
			messages := parseMessages(ircMessages)
			for _, m := range messages {
				if m.MessageType == AuthSuccess {
					auth <- nil
					return
				}
			}
		}
	}()
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_chat_auth_timeout)
		timeout <- errors.New("Failed to authenticate with Twitch chat")
	}()
	for {
		select {
		case err := <-auth:
			return err
		case err := <-timeout:
			return err
		}
	}
}

func (c *twitchChatConnection) isConnected() bool {
	return c.connection != nil || c.connecting
}

func (c *twitchChatConnection) disconnect() {
	c.sendQueue <- queueMessage{
		messageType: websocket.CloseMessage,
		message:     websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	}
}

func (c *twitchChatConnection) authenticate(token, username string) {
	c.send("CAP", "REQ", ":twitch.tv/membership twitch.tv/tags twitch.tv/commands")
	c.send("PASS", fmt.Sprintf("oauth:%s", token))
	c.send("NICK", username)
}

func (c *twitchChatConnection) joinChannel(channelName string) {
	if !c.isConnected() {
		c.connect()
	}
	c.send("JOIN", "#"+channelName)
}

func (c *twitchChatConnection) leaveChannel(channel string) {
	c.send("PART", "#"+channel)
}

func (c *twitchChatConnection) message(channel, message string) {
	c.send("PRIVMSG", "#"+channel, ":"+message)
}

func (c *twitchChatConnection) send(command string, params ...string) {
	message := []byte(command + " ")
	message = append(message, strings.Join(params, " ")...)
	c.sendQueue <- queueMessage{messageType: websocket.TextMessage, message: message}
}

func (c *twitchChatConnection) readMessages() {
	for {
		_, ircMessages, err := c.connection.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "close 1000") {
				break
			}
		}
		messages := parseMessages(ircMessages)
		for _, m := range messages {
			if m.MessageType == Ping {
				c.send("PONG", m.Params...)
			} else {
				c.messages <- m
			}
		}
	}
	c.connection = nil
}

func (c *twitchChatConnection) writeMessages() {
	for {
		message := <-c.sendQueue
		c.connection.WriteMessage(message.messageType, message.message)
	}
}

func parseMessages(bytes []byte) []TwitchMessage {
	messageStrings := strings.Split(string(bytes), "\r\n")
	messages := []TwitchMessage{}
	for _, s := range messageStrings {
		if s == "" {
			continue
		}
		messages = append(messages, NewChatMessage(s))
	}
	return messages
}
