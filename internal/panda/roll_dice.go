package panda

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func (bot *Bot) rollDice(command Command) string {
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
		return fmt.Sprintf("Chill, %s!  I can roll up to 10 dice with up to 100 sides.", command.Message.User.DisplayName)
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
		return fmt.Sprintf("P.A.N.D.A has rolled %d %d-sided dice for %s. They rolled %s = %d!", diceCount, diceSides, command.Message.User.DisplayName, joinedRolls, result)
	} else {
		if diceSides == 20 {
			return getD20Response(command, result)
		} else {
			return fmt.Sprintf("P.A.N.D.A has rolled a %d-sided dice for %s. They rolled a %d!", diceSides, command.Message.User.DisplayName, result)
		}
	}
}

func getD20Response(command Command, result int) string {
	switch result {
	case 20:
		return fmt.Sprintf("P.A.N.D.A has deemed %s worthy of a crit! GG!", command.Message.User.DisplayName)
	case 1:
		return fmt.Sprintf("P.A.N.D.A has deemed %s unworthy. Nat 1! Better make your sacrifice for P.A.N.D.A", command.Message.User.DisplayName)
	default:
		return fmt.Sprintf("P.A.N.D.A has rolled a D20 for %s. They rolled a %d!", command.Message.User.DisplayName, result)
	}
}
