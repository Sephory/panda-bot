package youtube

import "github.com/sephory/panda-bot/pkg/chat"

type YouTubeClientConfiguration struct {
	ClientId string
	Secret   string
}

var _ chat.ChatClient = &YouTubeClient{}

type YouTubeClient struct {
	config *YouTubeClientConfiguration
	api    *youTubeApi
}

func New(config *YouTubeClientConfiguration) *YouTubeClient {
	return &YouTubeClient{
		config: config,
		api:    newYouTubeApi(),
	}
}

// JoinChannel implements chat.ChatClient
func (c *YouTubeClient) JoinChannel(channelName string) chat.ChatChannel {
	channel := c.api.findChannel(channelName)
	c.api.findStreams(channel.Id)
	return NewYouTubeChannel(channelName, c.api)
}

// LeaveChannel implements chat.ChatClient
func (c *YouTubeClient) LeaveChannel(channelName string) {
	panic("unimplemented")
}

// GetName implements chat.ChatClient
func (c *YouTubeClient) GetName() string {
	return "YouTube"
}
