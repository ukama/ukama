package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SubscriberRepo interface {
	Add(subscriber *Subscriber) error
	Get(subscriberId string) (*Subscriber, error)
	Delete(subscriberId string) error
	Update(subscriberId string, sub Subscriber) (*Subscriber, error)
}

type subscriberRepo struct {
	Db sql.Db
}

func NewSubscriberRepo(db sql.Db) *subscriberRepo {
	return &subscriberRepo{
		Db: db,
	}
}

func (s *subscriberRepo) Add(pkg *Subscriber) error {
	db := s.Db.GetGormDb()
	result := db.Create(pkg)
	return result.Error
}

func (s *subscriberRepo) Get(subscriberId string) (*Subscriber, error) {
	var subscriber Subscriber
	err := s.Db.GetGormDb().Where("subscriber_id = ?", subscriberId).First(&subscriber).Error
	if err != nil {
		return nil, err
	}
	return &subscriber, nil

}

func (s *subscriberRepo) Delete(subscriberId string) error {
	result := s.Db.GetGormDb().Where("subscriber_id = ?", subscriberId).Delete(&Subscriber{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (b *subscriberRepo) Update(subscriberId string, sub Subscriber) (*Subscriber, error) {
	result := b.Db.GetGormDb().Where(sub).UpdateColumns(sub)
	if result.Error != nil {
		return nil, result.Error
	}

	return &sub, nil
}

