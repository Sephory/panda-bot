package database

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

const TABLE_SERVICES = "services"

var _ models.Model = &Service{}

type Service struct {
	models.BaseModel
	Name string
}

func (s *Service) TableName() string {
	return TABLE_SERVICES
}

func serviceQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Service{})
}

func (d *Database) FindServiceById(id string) *Service {
	service := &Service{}
	serviceQuery(d.dao).
		Where(dbx.HashExp{"id": id}).
		Limit(1).
		One(service)
	return service
}
