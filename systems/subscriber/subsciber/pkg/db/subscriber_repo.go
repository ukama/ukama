package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SubscriberRepo interface {
	Add(subscriber *Subscriber) error
	Get(subscriberId string) (*Subscriber, error)
	Delete(subscriberId string) error
}

type subscriberRepo struct {
	Db sql.Db
}

func NewSubscriberRepo(db sql.Db) *subscriberRepo {
	return &subscriberRepo{
		Db: db,
	}
}

func (r *subscriberRepo) Add(_package *Subscriber) error {
	result := r.Db.GetGormDb().Create(_package)

	return result.Error
}

func (s *subscriberRepo) Get(subscriberId string) (*Subscriber, error) {
	var subscriber Subscriber

	result := s.Db.GetGormDb().Where("subscriberId = ?", subscriberId).First(&subscriber)

	if result.Error != nil {
		return nil, result.Error
	}

	return &subscriber, nil
}



func (s *subscriberRepo) Delete(subscriberId string) error {
	result := s.Db.GetGormDb().Where("subscriberId = ?", subscriberId).Delete(&Subscriber{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

