package twitch

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const twitch_event_subscription_endpoint = "https://api.twitch.tv/helix/eventsub/subscriptions"
const twitch_users_endpoint = "https://api.twitch.tv/helix/users"

type twitchApi struct {
	clientId      string
	token         string
	httpClient    *http.Client
	subscriptions map[string][]*eventSubscription
}

func newTwitchApi(clientId string, token string) *twitchApi {
	return &twitchApi{
		clientId:      clientId,
		token:         token,
		httpClient:    http.DefaultClient,
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
		api.delete(url)
	}
}

func (api *twitchApi) getUserId(channelName string) *twitchUser {
	url := twitch_users_endpoint + "?login=" + channelName
	response, err := api.get(url)
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

	response, err := api.post(twitch_event_subscription_endpoint, request)
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

func (api *twitchApi) get(url string) (*http.Response, error) {
	request, err := api.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return api.httpClient.Do(request)
}

func (api *twitchApi) post(url string, body interface{}) (*http.Response, error) {
	request, err := api.newRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return api.httpClient.Do(request)
}

func (api *twitchApi) delete(url string) (*http.Response, error) {
	request, err := api.newRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return api.httpClient.Do(request)
}

func (api *twitchApi) newRequest(method string, url string, body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + api.token},
		"Client-Id":     {api.clientId},
	}
	return request, nil
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
