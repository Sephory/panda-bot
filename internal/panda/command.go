package panda

import (
	"strings"

	"github.com/sephory/panda-bot/pkg/chat"
)

type CommandType int

const (
	Unknown CommandType = iota
	HelloWorld
	RollDice
	Poll
	Vote
	Results
)

var command_type_map = map[string]CommandType{
	"hello":   HelloWorld,
	"roll":    RollDice,
	"poll":    Poll,
	"vote":    Vote,
	"results": Results,
}

type Command struct {
	CommandType CommandType
	CommandText string
	Params      []string
	Message       chat.Message
}

func NewCommand(event chat.Message) Command {
	commandText := strings.Split(event.Text[1:], " ")
	command := Command{
		CommandType: Unknown,
		CommandText: strings.ToLower(commandText[0]),
		Message:       event,
	}

	if val, ok := command_type_map[command.CommandText]; ok {
		command.CommandType = val
	}

	if len(commandText) > 1 {
		command.Params = commandText[1:]
	}
	return command
}
