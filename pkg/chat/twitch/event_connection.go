package twitch

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const twitch_event_host = "eventsub-beta.wss.twitch.tv:443"
const twitch_event_path = "ws"

type twitchEventConnection struct {
	connection *websocket.Conn
	events     chan interface{}
	sessionId  string
}

func (c *twitchEventConnection) connect() error {
	url := url.URL{Scheme: "wss", Host: twitch_event_host, Path: twitch_event_path}
	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if connection != nil {
		c.connection = connection
		c.events = make(chan interface{})
		go c.readEvents()
	}
	return err
}

func (c *twitchEventConnection) readEvents() {
	defer close(c.events)
	for {
		event := &twitchEvent{}
		err := c.connection.ReadJSON(event)
		if err != nil {
			log.Panic(err)
			break
		}
		if event.Metadata.MessageType == "session_keepalive" {
			continue
		} else {
			payload := event.getPayload()
			c.events <- payload
		}
	}
}
