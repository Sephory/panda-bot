package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

type PollResult struct {
	PollOptionId   string
	PollOptionText string
	Votes          int
}

func pollResultQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.DB().
		Select("poll_options.id poll_option_id", "poll_options.text poll_option_text", "COUNT(poll_votes.id) votes").
		From(TABLE_POLL_VOTES).
		InnerJoin(TABLE_POLL_OPTIONS, dbx.NewExp("poll_options.id = poll_votes.poll_option_id")).
		InnerJoin(TABLE_POLLS, dbx.NewExp("polls.id = poll_options.poll_id")).
		GroupBy("poll_options.id", "poll_options.text")
}

func (d *Database) GetResultsForPoll(pollId string) []*PollResult {
	results := []*PollResult{}
	pollResultQuery(d.dao).
		Where(dbx.HashExp{"polls.id": pollId}).
		All(&results)
	return results
}
