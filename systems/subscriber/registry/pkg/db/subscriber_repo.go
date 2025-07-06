/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type SubscriberRepo interface {
	Add(subscriber *Subscriber, nestedFunc func(*Subscriber, *gorm.DB) error) error
	Get(subscriberId uuid.UUID) (*Subscriber, error)
	GetByEmail(email string) (*Subscriber, error)
	Delete(subscriberId uuid.UUID) error
	Update(subscriberId uuid.UUID, sub Subscriber) error
	GetByNetwork(networkId uuid.UUID) ([]Subscriber, error)
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

func (s *subscriberRepo) Add(subscriber *Subscriber, nestedFunc func(*Subscriber, *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(subscriber, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(subscriber)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
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

func (s *subscriberRepo) GetByEmail(email string) (*Subscriber, error) {
	var subscriber Subscriber

	err := s.Db.GetGormDb().Where("email = ?", email).First(&subscriber).Error
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

func (b *subscriberRepo) Update(subscriberId uuid.UUID, sub Subscriber) error {
	db := b.Db.GetGormDb()
	err := db.Model(&Subscriber{}).Where("subscriber_id = ?", subscriberId).Updates(sub).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *subscriberRepo) GetByNetwork(networkId uuid.UUID) ([]Subscriber, error) {
	var subscribers []Subscriber
	result := s.Db.GetGormDb().Where(&Subscriber{NetworkId: networkId}).Find(&subscribers)

	if result.Error != nil {
		return nil, result.Error
	}
	return subscribers, nil
}
