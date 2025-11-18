/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type NetRepo interface {
	Add(network *Network, nestedFunc func(*Network, *gorm.DB) error) error
	Get(id uuid.UUID) (*Network, error)
	SetDefault(id uuid.UUID, isDefault bool) (*Network, error)
	GetByName(network string) (*Network, error)
	GetAll() ([]Network, error)
	GetDefault() (*Network, error)
	Delete(id uuid.UUID) error
	GetNetworkCount() (int64, error)
}

type netRepo struct {
	Db sql.Db
}

func NewNetRepo(db sql.Db) NetRepo {
	return &netRepo{
		Db: db,
	}
}

func (n netRepo) Get(id uuid.UUID) (*Network, error) {
	var ntwk Network

	result := n.Db.GetGormDb().First(&ntwk, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ntwk, nil
}

func (n netRepo) GetDefault() (*Network, error) {
	var ntwk Network

	result := n.Db.GetGormDb().Where("is_default = ?", true).First(&ntwk)
	if result.Error != nil {
		return nil, result.Error
	}

	return &ntwk, nil
}

func (n netRepo) SetDefault(id uuid.UUID, isDefault bool) (*Network, error) {
	var ntwk Network

	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&ntwk, id).Error; err != nil {
			return fmt.Errorf("failed to find network %v: %w", id, err)
		}

		if ntwk.IsDefault == isDefault {
			return nil
		}

		if isDefault {
			if err := tx.Model(&Network{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
				return fmt.Errorf("failed to clear existing default networks: %w", err)
			}
		}

		if err := tx.Model(&ntwk).Update("is_default", isDefault).Error; err != nil {
			return fmt.Errorf("failed to update network %v default status: %w", id, err)
		}

		ntwk.IsDefault = isDefault

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ntwk, nil
}

func (n netRepo) GetAll() ([]Network, error) {
	var ntwk []Network

	result := n.Db.GetGormDb().Model(&Network{}).Find(&ntwk)
	if result.Error != nil {
		return nil, result.Error
	}

	return ntwk, nil
}

func (n netRepo) GetByName(networkName string) (*Network, error) {
	var network Network

	result := n.Db.GetGormDb().Where("name = ?", networkName).First(&network)
	if result.Error != nil {
		return nil, result.Error
	}

	return &network, nil
}

func (n netRepo) Add(network *Network, nestedFunc func(*Network, *gorm.DB) error) error {
	return n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			if err := nestedFunc(network, tx); err != nil {
				return err
			}
		}

		return tx.Create(network).Error
	})
}

func (s netRepo) Delete(networkId uuid.UUID) error {
	result := s.Db.GetGormDb().Where("id = ?", networkId).Delete(&Network{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (n netRepo) GetNetworkCount() (int64, error) {
	var count int64
	result := n.Db.GetGormDb().Model(&Network{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
