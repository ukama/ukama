// This is an example of a repository
package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type ServiceRepo interface {
	Register(name string, url string) error
	UnRegister(serviceId string) error
	Update(serviceId string, url string) error
	Get(serviceId string) (*Service, error)
	List() (*Service, error)
}

type serviceRepo struct {
	db sql.Db
}

func NewServiceRepo(db sql.Db) *serviceRepo {
	return &serviceRepo{
		db: db,
	}
}

func (r *serviceRepo) Register(orgName string, url string) error {
	return nil
}

func (r *serviceRepo) Update(servieId string, url string) error {
	return nil
}

func (r *serviceRepo) UnRegister(serviceId string) error {
	return nil
}

func (r *serviceRepo) List() (*Service, error) {

	return nil, nil
}

func (r *serviceRepo) Get(serviceId string) (*Service, error) {
	return nil, nil
}
