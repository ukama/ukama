package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type SubscriberRepo interface {
	Add(subscriber *Subscriber) error
	Get(subscriberID uuid.UUID) (*Subscriber, error)
	Delete(subscriberID uuid.UUID) error
	Update(subscriberID uuid.UUID, sub Subscriber)  error
	GetByNetwork(networkID uuid.UUID) ([]Subscriber, error)
	ListSubscribers() ([]Subscriber, error)
}

type subscriberRepo struct {
	Db sql.Db
}

func NewSubscriberRepo(db sql.Db) SubscriberRepo {
	return &subscriberRepo{
		Db: db,
	}
}

func (s *subscriberRepo) Add(subscriber *Subscriber) error {
	db := s.Db.GetGormDb()
	return db.Create(subscriber).Error
}
func (s *subscriberRepo) ListSubscribers() ([]Subscriber, error) {

	var subscribers []Subscriber
	result := s.Db.GetGormDb().Find(&subscribers)
	if result.Error != nil {
		return nil, result.Error
	}
	return subscribers, nil
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

func (b *subscriberRepo) Update(subscriberID uuid.UUID, sub Subscriber)  error {
	var updatedSubscriber Subscriber
	db := b.Db.GetGormDb()
	err := db.Model(&Subscriber{}).Where("subscriber_id = ?", subscriberID).Updates(sub).First(&updatedSubscriber).Error
	if err != nil {
		return  err
	}
	return  nil
}

func (s *subscriberRepo) GetByNetwork(networkID uuid.UUID) ([]Subscriber, error) {
	var subscribers []Subscriber
	result := s.Db.GetGormDb().Where(&Subscriber{NetworkID: networkID}).Find(&subscribers)

	if result.Error != nil {
		return nil, result.Error
	}
	return subscribers, nil
}
