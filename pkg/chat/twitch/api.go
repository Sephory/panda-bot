package twitch

import (
	"encoding/json"
	"net/http"

	"github.com/sephory/panda-bot/pkg/chat"
)

const twitch_event_subscription_endpoint = "eventsub/subscriptions"
const twitch_users_endpoint = "users"

type twitchApi struct {
	client        *chat.ApiClient
	subscriptions map[string][]*eventSubscription
}

func newTwitchApi(clientId string, token string) *twitchApi {
	baseUrl := "https://api.twitch.tv/helix/"
	header := http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token},
		"Client-Id":     {clientId},
	}
	client, err := chat.NewApiClient(baseUrl, header)
	if err != nil {
		panic(err)
	}
	return &twitchApi{
		client:        client,
		subscriptions: map[string][]*eventSubscription{},
	}
}

func (api *twitchApi) subscribe(channelName string, sessionId string) {
	api.subscriptions[channelName] = []*eventSubscription{}
	if err := api.createSubscription(channelName, "channel.update", sessionId); err != nil {
		panic(err)
	}
}

func (api *twitchApi) unsubscribe(channelName string) {
	for _, subscription := range api.subscriptions[channelName] {
		url := twitch_event_subscription_endpoint + "?id=" + subscription.Id
		api.client.Delete(url, nil)
	}
}

func (api *twitchApi) getUserId(channelName string) *twitchUser {
	response, err := api.client.Get(
		twitch_users_endpoint,
		map[string]string{
			"login": channelName,
		},
	)
	if err != nil {
		panic(err)
	}
	user := &usersResponse{}
	err = json.NewDecoder(response.Body).Decode(user)
	if err != nil {
		panic(err)
	}
	return user.Data[0]
}

func (api *twitchApi) createSubscription(channelName string, subscriptionType string, sessionId string) error {
	user := api.getUserId(channelName)
	request := eventSubscriptionRequest{
		Type:    subscriptionType,
		Version: "1",
		Condition: broadcasterUserIdCondition{
			BroadcasterUserId: user.Id,
		},
		Transport: eventSubscriptionTransport{
			Method:    "websocket",
			SessionId: sessionId,
		},
	}

	response, err := api.client.Post(twitch_event_subscription_endpoint, nil, request)
	if err != nil {
		return err
	}
	subscribeResponse := &eventSubscriptionResponse{}
	err = json.NewDecoder(response.Body).Decode(subscribeResponse)
	if err != nil {
		return err
	}
	api.subscriptions[channelName] = append(api.subscriptions[channelName], subscribeResponse.Data[0])
	return nil
}

type twitchUser struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

type usersResponse struct {
	Data []*twitchUser `json:"data"`
}

type eventSubscriptionRequest struct {
	Type      string                     `json:"type"`
	Version   string                     `json:"version"`
	Condition interface{}                `json:"condition"`
	Transport eventSubscriptionTransport `json:"transport"`
}

type eventSubscription struct {
	Id        string                     `json:"id"`
	Status    string                     `json:"status"`
	Type      string                     `json:"type"`
	Version   string                     `json:"version"`
	Condition *json.RawMessage           `json:"condition"`
	CreatedAt string                     `json:"created_at"`
	Transport eventSubscriptionTransport `json:"transport"`
	Cost      int                        `json:"cost"`
}

type eventSubscriptionResponse struct {
	Data         []*eventSubscription `json:"data"`
	Total        int                  `json:"total"`
	TotalCost    int                  `json:"total_cost"`
	MaxTotalCost int                  `json:"max_total_cost"`
}

type eventSubscriptionTransport struct {
	Method      string `json:"method"`
	SessionId   string `json:"session_id"`
	ConnectedAt string `json:"connected_at"`
}

type broadcasterUserIdCondition struct {
	BroadcasterUserId string `json:"broadcaster_user_id"`
}
