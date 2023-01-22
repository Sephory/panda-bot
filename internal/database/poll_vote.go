package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_POLL_VOTES = "poll_votes"

var _ models.Model = &PollVote{}

type PollVote struct {
	models.BaseModel
	PollOptionId string
	ChannelId    string
	Voter        string
}

func (pv *PollVote) TableName() string {
	return TABLE_POLL_VOTES
}

func pollVoteQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&PollVote{})
}

func (d *Database) DidVoterVote(pollId string, channelId string, voter string) bool {
	rows, err := pollVoteQuery(d.dao).
		InnerJoin(TABLE_POLL_OPTIONS, dbx.NewExp("poll_options.id = poll_votes.poll_option_id")).
		Where(dbx.HashExp{
			"poll_votes.channel_id": channelId,
			"poll_votes.voter":      voter,
			"poll_options.poll_id":  pollId,
		}).
		Rows()
	if err != nil {
		return false
	}
	return rows.Next()
}

func (d *Database) SavePollVote(vote *PollVote) {
	d.dao.Save(vote)
}
