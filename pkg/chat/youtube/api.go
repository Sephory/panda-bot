package youtube

import (
	"context"
	"log"

	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

type youTubeApi struct {
	client *yt.Service
}

func newYouTubeApi() *youTubeApi {
	client, err := yt.NewService(context.Background(), option.WithCredentialsFile("youtube.json"))
	if err != nil {
		panic(err)
	}
	return &youTubeApi{
		client: client,
	}
}

func (api *youTubeApi) findChannel(channelName string) *yt.Channel {
	channels, err := api.client.Channels.List([]string{"id"}).
		ForUsername(channelName).
		Do()
	if err != nil {
		panic(err)
	}
	if len(channels.Items) == 0 {
		channels, err = api.client.Channels.List([]string{"id"}).
			Id(channelName).
			Do()
	}
	if len(channels.Items) == 0 {
		return nil
	}
	return channels.Items[0]
}

func (api *youTubeApi) getLiveChatIds(channelId string) []string {
	results, err := api.client.Search.List([]string{"snippet"}).
		ChannelId(channelId).
		Type("video").
		EventType("live").
		Do()
	if err != nil {
		panic(err)
	}
	var videos *yt.VideoListResponse
	for _, r := range results.Items {
		videos, err = api.client.Videos.List([]string{"liveStreamingDetails"}).
			Id(r.Id.VideoId).
			Do()
		if err != nil {
			panic(err)
		}
	}
	if videos == nil || len(videos.Items) == 0 {
		return nil
	}
	log.Printf("Found %v videos for %s", len(videos.Items), channelId)
	liveChatIds := []string{}
	for _, v := range videos.Items {
		liveChatIds = append(liveChatIds, v.LiveStreamingDetails.ActiveLiveChatId)
	}
	return liveChatIds
}

func (api *youTubeApi) getChat(liveChatId string, nextPageToken string) *yt.LiveChatMessageListResponse {
	messageRequest := api.client.LiveChatMessages.List(liveChatId, []string{"snippet", "authorDetails"})
	if nextPageToken != "" {
		messageRequest = messageRequest.PageToken(nextPageToken)
	}
	messages, err := messageRequest.Do()
	if err != nil {
		panic(err)
	}
	return messages
}

func (api *youTubeApi) sendChat(message string, liveChatId string) {

}
