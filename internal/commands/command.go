package commands

import (
	"strings"

	"github.com/sephory/panda-bot/pkg/chat"
)

type CommandType int

const (
	Unknown CommandType = iota
	HelloWorld
	RollDice
)

var command_type_map = map[string]CommandType{
	"hello":    HelloWorld,
	"roll": RollDice,
}

type Command struct {
	CommandType CommandType
	CommandText string
	Params      []string
	Event       chat.ChatEvent
}

func New(event chat.ChatEvent) Command {
	commandText := strings.Split(event.Message[1:], " ")
	command := Command{
		CommandType: Unknown,
		CommandText: strings.ToLower(commandText[0]),
		Event:       event,
	}

	if val, ok := command_type_map[command.CommandText]; ok {
		command.CommandType = val
	}

	if len(commandText) > 1 {
		command.Params = commandText[1:]
	}
	return command
}
