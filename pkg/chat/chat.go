package chat

type ChatClient interface {
	Connect() error
	Disconnect()
	JoinChannel(channelName string) ChatChannel
	LeaveChannel(channelName string)
	GetName() string
}

type ChatChannel interface {
	GetName() string
	GetEvents() chan interface{}
	SendMessage(message string)
}

type Message struct {
	User    UserInfo
	Message string
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
