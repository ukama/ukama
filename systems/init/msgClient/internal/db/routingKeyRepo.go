// This is an example of a repository
package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type RoutingKeyRepo interface {
	Add(key string) error
	Remove(key string) error
	Register(key string, serviceId string) error
	UnRegister(key string, serviceId string) error
	List() (*RoutingKey, error)
	ReadAllRoutes() ([]string, error)
	Get(key string) (*RoutingKey, error)
}

type routingKeyRepo struct {
	db sql.Db
}

func NewRoutingKeyRepo(db sql.Db) *routingKeyRepo {
	return &routingKeyRepo{
		db: db,
	}
}

func (r *routingKeyRepo) Add(key string) error {
	return nil
}

func (r *routingKeyRepo) Remove(key string) error {
	return nil
}

func (r *routingKeyRepo) Register(key string, serviceId string) error {
	return nil
}

func (r *routingKeyRepo) UnRegister(key string, serviceId string) error {
	return nil
}

func (r *routingKeyRepo) List() (*RoutingKey, error) {

	return nil, nil
}

func (r *routingKeyRepo) ReadAllRoutes() ([]string, error) {

	return nil, nil
}

func (r *routingKeyRepo) Get(key string) (*RoutingKey, error) {
	return nil, nil
}
