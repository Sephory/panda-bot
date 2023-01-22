package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_CUSTOM_RESPONSES = "custom_responses"

var _ models.Model = &CustomResponse{}

type CustomResponse struct {
	models.BaseModel
	UserId      string
	Command     string
	Response    string
	IsModOnly   bool
	IsOwnerOnly bool
}

func (cr *CustomResponse) TableName() string {
	return TABLE_CUSTOM_RESPONSES
}

func customResponseQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&CustomResponse{})
}

func (d *Database) GetCustomResponseForChannel(channelName string, commandText string) *CustomResponse {
	response := &CustomResponse{}
	customResponseQuery(d.dao).
		Join("INNER JOIN", TABLE_CHANNELS, dbx.NewExp("channels.user_id = custom_responses.user_id")).
		Where(dbx.HashExp{
			"channels.name":            channelName,
			"custom_responses.command": commandText,
		}).
		Limit(1).
		One(response)

	return response
}
