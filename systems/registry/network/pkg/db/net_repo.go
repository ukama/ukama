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
	"github.com/ukama/ukama/systems/common/validation"
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

	// Start a database transaction
	tx := n.Db.GetGormDb().Begin()

	// Set all networks to is_default false
	if err := tx.Model(&Network{}).Where("1 = 1").Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set all networks to not default: %w", err)
	}

	// Find the network with the id
	if err := tx.First(&ntwk, id).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find network %v: %w", id, err)
	}

	// Set the network to is_default true
	if err := tx.Model(&ntwk).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set network %v to default: %w", id, err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &network, nil
}

// func (n netRepo) GetByOrgName(orgID uint) ([]Network, error) {
// This gives the result in a single sql query, but fail to distingush between
// when org does not exist vs when org has no networks, can improve later.
// result := db.Joins("JOIN orgs on orgs.id=networks.org_id").
// Where("orgs.name=? and orgs.deleted_at is null", orgName).Debug().Find(&networks)
// }

func (n netRepo) Add(network *Network, nestedFunc func(network *Network, tx *gorm.DB) error) error {
	if !validation.IsValidDnsLabelName(network.Name) {
		return fmt.Errorf("invalid name. must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(network, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(network)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
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
