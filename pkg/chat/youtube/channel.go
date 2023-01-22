package youtube

import "github.com/sephory/panda-bot/pkg/chat"

var _ chat.ChatChannel = &YouTubeChannel{}

type YouTubeChannel struct {
	name   string
	api    *youTubeApi
	events chan interface{}
}

func NewYouTubeChannel(name string, api *youTubeApi) *YouTubeChannel {
	return &YouTubeChannel{
		name:   name,
		api:    api,
		events: make(chan interface{}),
	}
}

// GetEvents implements chat.ChatChannel
func (c *YouTubeChannel) GetEvents() chan interface{} {
	return c.events
}

// GetName implements chat.ChatChannel
func (c *YouTubeChannel) GetName() string {
	return c.name
}

// SendMessage implements chat.ChatChannel
func (c *YouTubeChannel) SendMessage(message string) {
	panic("unimplemented")
}
