package panda

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/sephory/panda-bot/internal/commands"
	"github.com/sephory/panda-bot/internal/data"
	"github.com/sephory/panda-bot/pkg/chat"
	"github.com/sephory/panda-bot/pkg/chat/twitch"
)

const (
	CLIENT_TWITCH = "Twitch"
)

type pandaBot struct {
	dao           *daos.Dao
	clients       map[string]chat.ChatClient
	commandPrefix byte
}

func newBot(config Config) *pandaBot {
	twitchClient := twitch.NewTwitchClient(config.Twitch)
	twitchClient.Connect()

	clients := map[string]chat.ChatClient{
		"Twitch": twitchClient,
	}
	return &pandaBot{
		clients:       clients,
		commandPrefix: config.CommandPrefix[0],
	}
}

func (bot *pandaBot) start(dao *daos.Dao) {
	bot.dao = dao
	bot.joinAll()
}

func (bot *pandaBot) joinAll() {
	var channels []data.Channel
	bot.dao.DB().NewQuery("SELECT * FROM channels").All(&channels)
	for _, channel := range channels {
		client := bot.clients[channel.Service]
		go bot.join(client, channel.Name)
	}
}

func (bot *pandaBot) join(client chat.ChatClient, channelName string) {
	channel := client.JoinChannel(channelName)
	log.Printf("Joining %s", channel.GetName())
	for event := range channel.GetEvents() {
		log.Printf("(%s) %s: %s", channel.GetName(), event.User.DisplayName, event.Message)
		if event.Message[0] == bot.commandPrefix {
			command := commands.New(event)
			bot.handleCommand(command, channel)
		}
	}
}

func (bot *pandaBot) handleCommand(command commands.Command, channel chat.ChatChannel) {
	switch command.CommandType {
	case commands.HelloWorld:
		channel.SendMessage(fmt.Sprintf("Hello, %s!", command.Event.User.DisplayName))
	case commands.RollDice:
		diceCount := 1
		diceSides := 20
		if len(command.Params) > 0 {
			diceMatch := regexp.MustCompile("\\d+d\\d+")
			multiDice := diceMatch.FindString(command.Params[0])
			if multiDice != "" {
				nums := strings.Split(multiDice, "d")
				diceCount, _ = strconv.Atoi(nums[0])
				diceSides, _ = strconv.Atoi(nums[1])
			}
			num, err := strconv.Atoi(command.Params[0])
			if err == nil && num > 0 {
				diceSides = num
			}
		}
		if diceCount > 10 || diceSides > 100 {
			channel.SendMessage(fmt.Sprintf("Chill, %s!  I can roll up to 10 dice with up to 100 sides.", command.Event.User.DisplayName))
		}
		rolls := []string{}
		result := 0
		for i := 0; i < diceCount; i++ {
			roll := rand.Intn(diceSides) + 1
			rolls = append(rolls, strconv.Itoa(roll))
			result += roll
		}
		if diceCount > 1 {
			joinedRolls := strings.Join(rolls, " + ")
			channel.SendMessage(fmt.Sprintf("P.A.N.D.A has rolled %d %d-sided dice for %s. They rolled %s = %d!", diceCount, diceSides, command.Event.User.DisplayName, joinedRolls, result))
		} else {
			if diceSides == 20 {
				channel.SendMessage(getD20Response(command, result))
			} else {
				channel.SendMessage(fmt.Sprintf("P.A.N.D.A has rolled a %d-sided dice for %s. They rolled a %d!", diceSides, command.Event.User.DisplayName, result))
			}
		}
	}
}

func getD20Response(command commands.Command, result int) string {
	switch result {
	case 20:
		return fmt.Sprintf("P.A.N.D.A has deemed %s worthy of a crit! GG!", command.Event.User.DisplayName)
	case 1:
		return fmt.Sprintf("P.A.N.D.A has deemed %s unworthy. Nat 1! Better make your sacrifice for P.A.N.D.A", command.Event.User.DisplayName)
	default:
		return fmt.Sprintf("P.A.N.D.A has rolled a D20 for %s. They rolled a %d!", command.Event.User.DisplayName, result)
	}
}
