package twitch

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

const twitch_chat_url = "irc-ws.chat.twitch.tv:443"

type twitchChatConnection struct {
	connection *websocket.Conn
	username   string
	messages   chan TwitchMessage
}

func (c *twitchChatConnection) connect() error {
	url := url.URL{Scheme: "wss", Host: twitch_chat_url}
	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if connection != nil {
		c.connection = connection
		c.messages = make(chan TwitchMessage)
		go c.readMessages()
	}
	return err
}

func (c *twitchChatConnection) disconnect() {
	c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *twitchChatConnection) authenticate(token, username string) {
	c.username = username
	c.send("CAP", "REQ", ":twitch.tv/membership twitch.tv/tags twitch.tv/commands")
	c.send("PASS", fmt.Sprintf("oauth:%s", token))
	c.send("NICK", username)
}

func (c *twitchChatConnection) joinChannel(channel string) {
	c.send("JOIN", "#"+channel)
}

func (c *twitchChatConnection) leaveChannel(channel string) {
	c.send("PART", "#"+channel)
}

func (c *twitchChatConnection) sendMessage(channel, message string) {
	c.send("PRIVMSG", "#"+channel, ":"+message)
}

func (c *twitchChatConnection) readMessages() {
	defer close(c.messages)
	for {
		_, ircMessages, err := c.connection.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "close 1000") {
				break
			}
		}
		messages := c.parseMessages(ircMessages)
		for _, m := range messages {
			if m.MessageType == Ping {
				c.send("PONG", m.Params...)
			} else {
				c.messages <- m
			}
		}
	}
}

func (c *twitchChatConnection) parseMessages(bytes []byte) []TwitchMessage {
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

func (c *twitchChatConnection) send(command string, params ...string) error {
	message := []byte(command + " ")
	message = append(message, strings.Join(params, " ")...)
	err := c.connection.WriteMessage(websocket.TextMessage, message)
	return err
}
