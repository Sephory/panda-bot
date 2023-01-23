package twitch

import (
	"errors"
	"strings"

	"github.com/sephory/panda-bot/pkg/chat"
)

type TwitchMessageType int

const (
	Unknown TwitchMessageType = iota
	Capabilities
	AuthSuccess
	GlobalUserState
	Join
	UserState
	RoomState
	ChatMessage
	Ping
	Notice
)

var message_type_map = map[string]TwitchMessageType{
	"CAP":             Capabilities,
	"001":             AuthSuccess,
	"GLOBALUSERSTATE": GlobalUserState,
	"JOIN":            Join,
	"USERSTATE":       UserState,
	"ROOMSTATE":       RoomState,
	"PRIVMSG":         ChatMessage,
	"PING":            Ping,
	"NOTICE":          Notice,
}

type MessageSource struct {
	FullName string
	Name     string
}

type TwitchMessage struct {
	MessageType TwitchMessageType
	Command     string
	Params      []string
	Source      MessageSource
	Tags        map[string][]string
}

func NewChatMessage(messageText string) TwitchMessage {
	message := TwitchMessage{
		MessageType: Unknown,
	}

	messageParts := strings.SplitN(messageText, " ", 2)

	if messageParts[0][0] == '@' {
		message.Tags = parseMessageTags(messageParts[0])
		messageParts = strings.SplitN(messageParts[1], " ", 2)
	}

	if messageParts[0][0] == ':' {
		message.Source = parseMessageSource(messageParts[0])
		messageParts = strings.SplitN(messageParts[1], " ", 2)
	}

	message.Command = messageParts[0]
	message.MessageType = parseMessageType(messageParts[0])

	if len(messageParts) > 1 {
		paramsParts := strings.SplitN(messageParts[1], ":", 2)
		message.Params = strings.Split(strings.TrimRight(paramsParts[0], " "), " ")
		if len(paramsParts) > 1 {
			message.Params = append(message.Params, paramsParts[1])
		}
	}
	return message
}

func (m *TwitchMessage) toChatEvent() (interface{}, error) {
	switch m.MessageType {
	case ChatMessage:
		return chat.Message{
			User:    m.toUserInfo(),
			Text: m.Params[1],
		}, nil
	case Join:
		return chat.UserJoin{
			User: m.toUserInfo(),
		}, nil
	default:
		return chat.Message{}, errors.New("Can't convert to ChatEvent")
	}
}

func (m *TwitchMessage) toUserInfo() chat.UserInfo {
	userInfo := chat.UserInfo{
		Username: m.Source.Name,
	}
	if m.MessageType == ChatMessage {
		userInfo.DisplayName = m.Tags["display-name"][0]
		userInfo.IsMod = m.Tags["mod"][0] == "1"
		userInfo.IsSubscriber = m.Tags["subscriber"][0] == "1"
	}
	return userInfo
}

func parseMessageTags(tagsText string) map[string][]string {
	tags := map[string][]string{}
	allTags := strings.Split(tagsText[1:], ";")
	for _, t := range allTags {
		keyValue := strings.Split(t, "=")
		tags[keyValue[0]] = strings.Split(keyValue[1], ",")
	}
	return tags
}

func parseMessageSource(sourceText string) MessageSource {
	split := strings.Index(sourceText, "!")
	if split > 0 {
		return MessageSource{
			FullName: sourceText[split+1:],
			Name:     sourceText[1:split],
		}
	} else {
		return MessageSource{
			FullName: sourceText,
		}
	}
}

func parseMessageType(typeText string) TwitchMessageType {
	if val, ok := message_type_map[typeText]; ok {
		return val
	}
	return Unknown
}
