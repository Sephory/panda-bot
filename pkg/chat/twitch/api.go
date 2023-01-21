package twitch

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const twitch_event_subcribe_endpoint = "https://api.twitch.tv/helix/eventsub/subscriptions"
const twitch_users_endpoint = "https://api.twitch.tv/helix/users"

type twitchApi struct {
	config     *TwitchClientConfiguration
	httpClient *http.Client
}

func (api *twitchApi) subscribe(channelName string, sessionId string) {
	user := api.getUserId(channelName)
	if err := api.createSubscription("channel.update", user.Id, sessionId); err != nil {
		panic(err)
	}
}

func (api *twitchApi) unsubscribe(channelName string) {

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

func (api *twitchApi) createSubscription(subscriptionType string, broadcasterUserId string, sessionId string) error {
	request := eventSubscriptionRequest{
		Type:    subscriptionType,
		Version: "1",
		Condition: broadcasterUserIdCondition{
			BroadcasterUserId: broadcasterUserId,
		},
		Transport: eventSubscriptionTransport{
			Method:    "websocket",
			SessionId: sessionId,
		},
	}

	api.post(twitch_event_subcribe_endpoint, request)
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

func (api *twitchApi) newRequest(method string, url string, body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+api.config.Token)
	request.Header.Add("Client-Id", api.config.ClientId)
	return request, nil
}

type twitchUser struct {
	Id              string
	Login           string
	DisplayName     string `json:"display_name"`
	Type            string
	BroadcasterType string `json:"broadcaster_type"`
	Description     string
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string
	CreatedAt       string `json:"created_at"`
}

type usersResponse struct {
	Data []*twitchUser
}

type eventSubscriptionRequest struct {
	Type      string                     `json:"type"`
	Version   string                     `json:"version"`
	Condition interface{}                `json:"condition"`
	Transport eventSubscriptionTransport `json:"transport"`
}

type eventSubscriptionResponse struct {
	Data struct {
		Id        string
		Status    string
		Type      string
		Version   string
		Condition *json.RawMessage
		CreatedAt string `json:"created_at"`
		Transport eventSubscriptionTransport
		Cost      int
	}
	Total        int
	TotalCost    int `json:"total_cost"`
	MaxTotalCost int `json:"max_total_cost"`
}

type eventSubscriptionTransport struct {
	Method      string `json:"method"`
	SessionId   string `json:"session_id"`
	ConnectedAt string `json:"connected_at"`
}

type broadcasterUserIdCondition struct {
	BroadcasterUserId string `json:"broadcaster_user_id"`
}
