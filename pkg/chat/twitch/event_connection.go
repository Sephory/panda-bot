package twitch

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const twitch_event_host = "eventsub-beta.wss.twitch.tv:443"
const twitch_event_path = "ws"
const twitch_event_welcome_timeout = 500

type twitchEventConnection struct {
	connection *websocket.Conn
	connecting bool
	events     chan interface{}
	sessionId  string
}

func newTwitchEventConnection() *twitchEventConnection {
	return &twitchEventConnection{
		events: make(chan interface{}),
	}
}

func (c *twitchEventConnection) connect() error {
	c.connecting = true
	defer func() { c.connecting = false }()
	url := url.URL{Scheme: "wss", Host: twitch_event_host, Path: twitch_event_path}
	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if connection != nil {
		c.connection = connection
		err = c.awaitEventSession()
		if err != nil {
			return err
		}
		go c.readEvents()
	}
	return err
}

func (c *twitchEventConnection) awaitEventSession() error {
	session := make(chan error, 1)
	go func() {
		for {
			event := &twitchEvent{}
			err := c.connection.ReadJSON(event)
			if err != nil {
				session <- err
				break
			}
			payload := event.getPayload()
			if welcome, ok := payload.(*sessionWelcome); ok {
				c.sessionId = welcome.Session.Id
				session <- nil
				break
			}
		}
	}()
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_event_welcome_timeout)
		timeout <- errors.New("Twitch EventSub welcome timed out")
	}()
	for {
		select {
		case err := <-session:
			return err
		case err := <-timeout:
			return err
		}
	}
}

func (c *twitchEventConnection) disconnect() {
	c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *twitchEventConnection) isConnected() bool {
	return c.connection != nil
}

func (c *twitchEventConnection) getSessionId() string {
	for !c.isConnected() {
		if (!c.connecting) {
			c.connect()
		} else {
			time.Sleep(time.Millisecond * 500)
		}
	}
	return c.sessionId
}

func (c *twitchEventConnection) readEvents() {
	for {
		event := &twitchEvent{}
		err := c.connection.ReadJSON(event)
		if err != nil {
			if strings.Contains(err.Error(), "close 1000") {
				break
			}
		}
		if event.Metadata.MessageType == "session_keepalive" {
			continue
		} else {
			payload := event.getPayload()
			c.events <- payload
		}
	}
	c.connection = nil
}
