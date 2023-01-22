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
	return channels.Items[0]
}

func (api *youTubeApi) findStreams(channelId string) {
	results, err := api.client.Search.List([]string{"snippet"}).
		ChannelId(channelId).
		Type("video").
		EventType("live").
		Do()
	if err != nil {
		panic(err)
	}
	log.Printf("Found %v broadcasts for %s", len(results.Items), channelId)
	for _, r := range results.Items {
		log.Printf("Live VideoId: %s", r.Id.VideoId)
	}
}
