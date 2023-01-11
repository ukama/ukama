// This is an example of a repository
package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type RouteRepo interface {
	Add(key string) (*Route, error)
	Remove(key string) error
	List() ([]Route, error)
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

func (r *routeRepo) Add(key string) (*Route, error) {

	var route Route
	res := r.db.GetGormDb().FirstOrCreate(&route, Route{Key: key})
	if res.Error != nil {
		return nil, res.Error
	}

	return &route, nil

}

func (r *routeRepo) Remove(key string) error {
	var rt Route
	res := r.db.GetGormDb().Delete(&rt, Route{Key: key})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *routeRepo) List() ([]Route, error) {

	var rt []Route
	res := r.db.GetGormDb().Find(&rt)
	if res.Error != nil {
		return nil, res.Error
	}

	return rt, nil
}

func (r *routeRepo) Get(key string) (*Route, error) {

	rt := Route{
		Key: key,
	}

	res := r.db.GetGormDb().Find(&rt)
	if res.Error != nil {
		return nil, res.Error
	}

	return &rt, nil
}
