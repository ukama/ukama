// This is an example of a repository
package db

import (
	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm/clause"
)

// declare interface so that we can mock it
type ServiceRepo interface {
	Register(service *Service) (*Service, error)
	UnRegister(serviceId string) error
	Get(serviceUuid string) (*Service, error)
	GetRoutes(serviceId string) ([]Route, error)
	List() ([]Service, error)
	AddRoute(s *Service, rt *Route) error
	RemoveRoutes(service *Service) error
}

type serviceRepo struct {
	db sql.Db
}

func NewServiceRepo(db sql.Db) *serviceRepo {
	return &serviceRepo{
		db: db,
	}
}

func (r *serviceRepo) Register(service *Service) (*Service, error) {

	res := r.db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"instance_id", "msg_bus_uri", "list_queue", "publ_queue", "exchange", "service_uri", "grpc_timeout", "deleted_at"}),
	}).Create(service)
	if res.Error != nil {

		return nil, res.Error
	}

	return service, nil
}

func (r *serviceRepo) UnRegister(serviceUuid string) error {
	var svc Service
	res := r.db.GetGormDb().Delete(&svc, Service{ServiceUuid: serviceUuid})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *serviceRepo) List() ([]Service, error) {
	var svc []Service
	res := r.db.GetGormDb().Preload("Routes").Find(&svc)
	if res.Error != nil {
		return nil, res.Error
	}

	return svc, nil
}

func (r *serviceRepo) GetRoutes(serviceUuid string) ([]Route, error) {
	var serviceRoutes []Route
	err := r.db.GetGormDb().Model(&Service{}).Where("service_uuid = ?", serviceUuid).Association("Routes").Find(&serviceRoutes)
	if err != nil {
		return nil, err
	}

	return serviceRoutes, nil
}

func (r *serviceRepo) Get(serviceUuid string) (*Service, error) {
	var svc Service
	res := r.db.GetGormDb().Preload("Routes").Where("service_uuid = ?", serviceUuid).First(&svc)
	if res.Error != nil {
		return nil, res.Error
	}

	return &svc, nil
}

func (r *serviceRepo) AddRoute(s *Service, rt *Route) error {

	res := r.db.GetGormDb().Model(s).Association("Routes").Append(rt)
	return res
}

func (r *serviceRepo) RemoveRoutes(service *Service) error {

	var serviceRoutes []Route
	err := r.db.GetGormDb().Model(service).Where("service_id = ?", service.ID).Association("Routes").Find(&serviceRoutes)
	if err != nil {
		return err
	}

	log.Infof("found old routes %+v", serviceRoutes)
	err = r.db.GetGormDb().Model(service).Association("Routes").Delete(&serviceRoutes)
	if err != nil {
		return err
	}

	return nil
}
