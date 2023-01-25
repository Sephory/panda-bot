package chat

type ChatClient interface {
	GetName() string
	JoinChannel(channelName string) ChatChannel
	LeaveChannel(channelName string)
}

type ChatChannel interface {
	GetName() string
	GetEvents() chan interface{}
	SendMessage(message string, options interface{})
}

type Message struct {
	User    UserInfo
	Text string
	Options interface{}
}

type UserJoin struct {
	User UserInfo
}

type UserInfo struct {
	Username     string
	DisplayName  string
	IsMod        bool
	IsSubscriber bool
}
