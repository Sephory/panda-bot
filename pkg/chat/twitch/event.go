package twitch

import "encoding/json"

type twitchEvent struct {
	Metadata struct {
		MessageId        string `json:"message_id"`
		MessageType      string `json:"message_type"`
		MessageTimestamp string `json:"message_timestamp"`
	}
	Payload *json.RawMessage
}

var payloadHandlers = map[string]func() interface{}{
	"session_welcome": func() interface{} { return &sessionWelcome{} },
}

func (e *twitchEvent) getPayload() interface{} {
	if e.Payload == nil {
		return nil
	}
	handler := payloadHandlers[e.Metadata.MessageType]
	if handler == nil {
		return nil
	}
	payload := handler()
	json.Unmarshal(*e.Payload, payload)
	return payload
}

type sessionWelcome struct {
	Session struct {
		Id                      string
		Status                  string
		ConnectedAt             string `json:"connected_at"`
		KeepaliveTimeoutSeconds int    `json:"keepalive_timeout_seconds"`
	}
}
