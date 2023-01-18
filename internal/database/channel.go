package database

import "github.com/pocketbase/pocketbase/models"

type Channel struct {
	models.BaseModel
	UserId string
	Service string
	Name string
}

func (c *Channel) TableName() string {
	return "channels"
}
