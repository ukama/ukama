// This is an example of a repository
package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

// declare interface so that we can mock it
type ServiceRepo interface {
	Register(service *Service) error
	UnRegister(serviceId string) error
	Update(serviceId string, url string) error
	Get(serviceId string) (*Service, error)
	GetRoutes(serviceId string) ([]Route, error)
	List() ([]Service, error)
}

type serviceRepo struct {
	db sql.Db
}

func NewServiceRepo(db sql.Db) *serviceRepo {
	return &serviceRepo{
		db: db,
	}
}

func (r *serviceRepo) Register(service *Service) error {

	res := r.db.GetGormDb().Create(service)
	if res.Error != nil {

		return res.Error
	}

	return nil
}

func (r *serviceRepo) Update(servieId string, url string) error {
	return nil
}

func (r *serviceRepo) UnRegister(serviceId string) error {
	return nil
}

func (r *serviceRepo) List() ([]Service, error) {

	return nil, nil
}

func (r *serviceRepo) GetRoutes(serviceId string) ([]Route, error) {
	var serviceRoutes []Route
	err := r.db.GetGormDb().Model(&Service{}).Where("service_id = ?", serviceId).Association("Routes").Find(&serviceRoutes)
	if err != nil {
		return nil, err
	}

	return serviceRoutes, nil
}

func (r *serviceRepo) Get(serviceId string) (*Service, error) {

	return nil, nil
}

func (r *serviceRepo) AddRoute(serviceId string, key string) error {

	service := Service{
		ServiceId: serviceId,
	}

	route := Route{
		Key: key,
	}

	res := r.db.GetGormDb().Model(&service).Association("Routes").Append(&route)
	return res
}

func (r *serviceRepo) RemoveRoutes(serviceId string, key string) error {

	service := Service{
		ServiceId: serviceId,
	}

	var routes []Route

	res := r.db.GetGormDb().Model(&Route{}).Where("service_id = ?", serviceId).Find(&routes)
	if res.Error != nil {
		return res.Error
	}

	err := r.db.GetGormDb().Model(&service).Association("Routes").Delete(&routes)
	if err != nil {
		return err
	}

	return nil
}
