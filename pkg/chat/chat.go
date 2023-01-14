package chat

type ChatClient interface {
	Connect() error
	Disconnect()
	JoinChannel(channelName string) ChatChannel
	LeaveChannel(channelName string)
}

type ChatChannel interface {
	GetName() string
	GetEvents() chan ChatEvent
	SendMessage(message string)
}

type ChatEvent struct {
	User UserInfo
	Message string
}

type UserInfo struct {
	Username     string
	DisplayName  string
	IsMod        bool
	IsSubscriber bool
}

