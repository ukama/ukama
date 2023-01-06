// This is an example of a repository
package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type RouteRepo interface {
	Add(key string) error
	Remove(key string) error
	Register(key string, serviceId string) error
	UnRegister(key string, serviceId string) error
	List() (*Route, error)
	Get(key string) (*Route, error)
}

type routeRepo struct {
	db sql.Db
}

func NewRouteRepo(db sql.Db) *routeRepo {
	return &routeRepo{
		db: db,
	}
}

func (r *routeRepo) Add(key string) error {
	var rt Route
	res := r.db.GetGormDb().FirstOrCreate(&rt, Route{Key: key})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *routeRepo) Remove(key string) error {
	return nil
}

func (r *routeRepo) Register(key string, serviceId string) error {
	return nil
}

func (r *routeRepo) UnRegister(key string, serviceId string) error {
	return nil
}

func (r *routeRepo) List() (*Route, error) {

	return nil, nil
}

func (r *routeRepo) Get(key string) (*Route, error) {
	return nil, nil
}
