package panda

import (
	"fmt"
	"strings"

	"github.com/sephory/panda-bot/internal/database"
)

func (bot *Bot) poll(command Command, channel *database.Channel) string {
	if len(command.Params) == 0 {
		return bot.activePolls(channel)
	}
	switch command.Params[0] {
	case "start":
		return bot.pollStart(command, channel)
	case "end":
		return bot.pollEnd(command, channel)
	default:
		return bot.pollInfo(command, channel)
	}
}

func (bot *Bot) activePolls(channel *database.Channel) string {
	polls := bot.db.GetActivePollsForUser(channel.UserId)
	switch len(polls) {
	case 0:
		return "There are no active polls"
	case 1:
		return fmt.Sprintf("%s  Type '!vote' followed by your answer to place your vote", polls[0].Prompt)
	default:
		pollNames := []string{}
		for _, p := range polls {
			pollNames = append(pollNames, "'"+p.Name+"'")
		}
		return fmt.Sprintf("There are a few polls to vote on. For more info, type '!poll' followed by %s", strings.Join(pollNames, " or "))
	}
}

func (bot *Bot) pollInfo(command Command, channel *database.Channel) string {
	poll := bot.db.GetPollForUserByName(channel.UserId, command.Params[0], true)
	if poll.IsActive {
		return fmt.Sprintf("%s  Type '!vote %s' followed by your answer to place your vote", poll.Prompt, poll.Name)
	}
	return ""
}

func (bot *Bot) pollStart(command Command, channel *database.Channel) string {
	if len(command.Params) < 3 {
		return ""
	}
	poll := bot.db.GetPollForUserByName(channel.UserId, command.Params[1], true)
	if poll.IsActive {
		return fmt.Sprintf("There is already and active poll called '%s'", poll.Name)
	}
	poll = &database.Poll{
		UserId:   channel.UserId,
		Name:     command.Params[1],
		Prompt:   strings.Join(command.Params[2:], " "),
		IsActive: true,
		IsOpen:   true,
	}
	bot.db.SavePoll(poll)
	return fmt.Sprintf("Poll started: %s", poll.Prompt)
}

func (bot *Bot) pollEnd(command Command, channel *database.Channel) string {
	var poll *database.Poll
	if len(command.Params) < 2 {
		poll = bot.db.GetOldestActivePollForUser(channel.UserId)
		if !poll.IsActive {
			return "I couldn't find an active poll"
		}
	} else {
		poll = bot.db.GetPollForUserByName(channel.UserId, command.Params[1], true)
		if !poll.IsActive {
			return fmt.Sprintf("I couldn't find an active poll called '%s'", command.Params[1])
		}
	}

	poll.IsActive = false
	bot.db.SavePoll(poll)
	return fmt.Sprintf("Poll ended: %s", poll.Prompt)
}

func (bot *Bot) vote(command Command, channel *database.Channel) string {
	polls := bot.db.GetActivePollsForUser(channel.UserId)
	if len(polls) == 0 {
		return "There's nothing to vote for at this time"
	}
	if len(command.Params) == 0 {
		return bot.activePolls(channel)
	}
	var poll *database.Poll
	var voteText string
	if len(polls) > 1 {
		if len(command.Params) < 2 {
			return bot.activePolls(channel)
		}
		name := command.Params[0]
		voteText = strings.Join(command.Params[1:], " ")
		for _, p := range polls {
			if p.Name == name {
				poll = p
				break
			}
		}
		if poll == nil {
			return bot.activePolls(channel)
		}
	} else {
		voteText = strings.Join(command.Params, " ")
		poll = polls[0]
	}
	voter := command.Event.User.Username

	if bot.db.DidVoterVote(poll.Id, channel.Id, voter) {
		return fmt.Sprintf("You already voted, %s!", command.Event.User.DisplayName)
	}

	options := bot.db.GetPollOptionsForPoll(poll.Id)
	var option *database.PollOption
	for _, o := range options {
		if strings.ToLower(strings.ReplaceAll(o.Text, " ", "")) == strings.ToLower(strings.ReplaceAll(voteText, " ", "")) {
			option = o
		}
	}

	if option == nil {
		if poll.IsOpen {
			option = &database.PollOption{
				PollId:    poll.Id,
				Text:      voteText,
				ChannelId: channel.Id,
				CreatedBy: voter,
			}
			bot.db.SavePollOption(option)

		} else {
			return fmt.Sprintf("'%s' is not a valid option for this poll", voteText)
		}
	}

	vote := &database.PollVote{
		PollOptionId: option.Id,
		ChannelId:    channel.Id,
		Voter:        voter,
	}

	bot.db.SavePollVote(vote)

	return fmt.Sprintf("Your vote has been counted, %s!", command.Event.User.DisplayName)
}

func (bot *Bot) results(command Command, channel *database.Channel) string {
	var poll *database.Poll
	if len(command.Params) > 0 {
		poll = bot.db.GetPollForUserByName(channel.UserId, command.Params[0], false)
	}
	if poll == nil || poll.Name == "" {
		poll = bot.db.GetMostRecentPollForUser(channel.UserId)
	}
	if poll.Name == "" {
		return "There are no polls to get results for!"
	}
	results := bot.db.GetResultsForPoll(poll.Id)
	if len(results) == 0 {
		return "There are no results to show!"
	}
	resultText := []string{}
	for _, r := range results {
		resultText = append(resultText, fmt.Sprintf("%s: %v", r.PollOptionText, r.Votes))
	}

	return fmt.Sprintf("Results for: %s  %s", poll.Prompt, strings.Join(resultText, ", "))
}
