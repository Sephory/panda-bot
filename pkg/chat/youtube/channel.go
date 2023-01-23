package youtube

import (
	"time"

	"github.com/sephory/panda-bot/pkg/chat"
	yt "google.golang.org/api/youtube/v3"
)

var _ chat.ChatChannel = &YouTubeChatChannel{}

type YouTubeMessageOptions struct {
	LiveChatId string
}

type YouTubeChatChannel struct {
	name   string
	channelId string
	api    *youTubeApi
	events chan interface{}
}

func NewYouTubeChannel(name string, channelId string, api *youTubeApi) *YouTubeChatChannel {
	channel := &YouTubeChatChannel{
		name:   name,
		channelId: channelId,
		api:    api,
		events: make(chan interface{}),
	}
	go channel.getLiveChats()
	return channel
}

// GetEvents implements chat.ChatChannel
func (c *YouTubeChatChannel) GetEvents() chan interface{} {
	return c.events
}

// GetName implements chat.ChatChannel
func (c *YouTubeChatChannel) GetName() string {
	return c.name
}

// SendMessage implements chat.ChatChannel
func (c *YouTubeChatChannel) SendMessage(message string, options interface{}) {
	if youTubeMessageOptions, ok := options.(YouTubeMessageOptions); ok {
		c.api.sendChat(message, youTubeMessageOptions.LiveChatId)
	}
}

func (c *YouTubeChatChannel) getLiveChats() {
	liveChatIds := c.api.getLiveChatIds(c.channelId)
	for _, id := range liveChatIds {
		go c.readMessages(id)
	}
}

func (c *YouTubeChatChannel) readMessages(liveChatId string) {
	response := c.api.getChat(liveChatId, "")
	nextPageToken := response.NextPageToken
	waitMs := response.PollingIntervalMillis
	for {
		time.Sleep(time.Millisecond * time.Duration(waitMs))
		response := c.api.getChat(liveChatId, nextPageToken)
		nextPageToken = response.NextPageToken
		waitMs = response.PollingIntervalMillis

		messages := response.Items
		chatEnded := false
		for _, m := range messages {
			if m.Snippet.Type == "chatEndedEvent" {
				chatEnded = true
			}
			c.events <- liveChatMessageToChatEvent(m)
		}
		if chatEnded {
			break
		}
	}
}

func liveChatMessageToChatEvent(liveChatMessage *yt.LiveChatMessage) interface{} {
	switch liveChatMessage.Snippet.Type {
	case "textMessageEvent":
		return chat.Message{
			User: authorDetailsToUserInfo(liveChatMessage.AuthorDetails),
			Text: liveChatMessage.Snippet.DisplayMessage,
			Options: YouTubeMessageOptions{
				LiveChatId: liveChatMessage.Snippet.LiveChatId,
			},
		}
	default:
		return nil
	}
}

func authorDetailsToUserInfo(authorDetails *yt.LiveChatMessageAuthorDetails) chat.UserInfo{
	return chat.UserInfo{
		Username: authorDetails.ChannelId,
		DisplayName: authorDetails.DisplayName,
		IsMod: authorDetails.IsChatModerator,
		IsSubscriber: authorDetails.IsChatSponsor,
	}
}
