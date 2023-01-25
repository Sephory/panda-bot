package twitch

import (
	"encoding/json"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/sephory/panda-bot/pkg/chat"
)

const twitch_event_url = "wdd://eventsub-beta.wss.twitch.tv:443/ws"
const twitch_event_welcome_timeout = 500

type twitchEventConnection struct {
	connection *chat.WebsocketConnection
	connecting sync.Mutex
	events     map[string]chan interface{}
	sessionId  string
}

func newTwitchEventConnection() *twitchEventConnection {
	eventsUrl, _ := url.Parse(twitch_event_url)
	return &twitchEventConnection{
		connection: chat.NewWebsocketConnection(eventsUrl),
		events:     make(map[string]chan interface{}),
	}
}

func (c *twitchEventConnection) awaitEventSession() error {
	c.connecting.Lock()
	defer c.connecting.Unlock()
	if c.sessionId != "" {
		return nil
	}
	timeout := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * twitch_event_welcome_timeout)
		timeout <- errors.New("Twitch EventSub welcome timed out")
	}()
	select {
	case sessionId := <-c.getSessionId():
		c.sessionId = sessionId
		go c.readEvents()
		return nil
	case err := <-timeout:
		return err
	}
}

func (c *twitchEventConnection) getSessionId() chan string {
	result := make(chan string, 1)
	go func() {
		messages, err := c.connection.GetMessages()
		if err != nil {
			result <- ""
		}
		for message := range messages {
			event := &twitchEvent{}
			err = json.Unmarshal(message.Bytes, event)
			if err != nil {
				result <- ""
				return
			}
			payload := event.getPayload()
			if welcome, ok := payload.(*sessionWelcome); ok {
				result <- welcome.Session.Id
				break
			}
		}
	}()
	return result
}

func (c *twitchEventConnection) readEvents() {
	messages, err := c.connection.GetMessages()
	if err != nil {
		return
	}
	for message := range messages {
		event := &twitchEvent{}
		json.Unmarshal(message.Bytes, event)
		if event.Metadata.MessageType == "session_keepalive" {
			continue
		} else {
		}
	}
	c.connection = nil
}
