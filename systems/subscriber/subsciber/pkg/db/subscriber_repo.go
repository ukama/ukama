package db

import (
	"github.com/gofrs/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SubscriberRepo interface {
	Add(subscriber *Subscriber) error
	Get(subscriberID uuid.UUID) (*Subscriber, error)
	Delete(subscriberID uuid.UUID) error
	Update(subscriberID uuid.UUID, sub Subscriber) (uuid.UUID, error)
	GetByNetwork(networkID uuid.UUID) ([]Subscriber, error)
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

func (s *subscriberRepo) Get(subscriberId uuid.UUID) (*Subscriber, error) {
	var subscriber Subscriber
	err := s.Db.GetGormDb().Where("subscriber_id = ?", subscriberId).First(&subscriber).Error
	if err != nil {
		return nil, err
	}
	return &subscriber, nil

}

func (s *subscriberRepo) Delete(subscriberId uuid.UUID) error {
	result := s.Db.GetGormDb().Where("subscriber_id = ?", subscriberId).Delete(&Subscriber{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (b *subscriberRepo) Update(subscriberId uuid.UUID, sub Subscriber) (uuid.UUID, error) {

	result := b.Db.GetGormDb().Where("subscriber_id = ?", subscriberId).UpdateColumns(sub)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return subscriberId, nil
}
func (s *subscriberRepo) GetByNetwork(networkID uuid.UUID) ([]Subscriber, error) {
	var subscribers []Subscriber
	result := s.Db.GetGormDb().Where(&Subscriber{NetworkID: networkID}).Find(&subscribers)

	if result.Error != nil {
		return nil, result.Error
	}
	return subscribers, nil
}
