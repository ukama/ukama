package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"github.com/ukama/ukama/systems/registry/network/pkg"
	"gorm.io/gorm"
)

type NetRepo interface {
	Add(network *Network, nestedFunc func(*Network, *gorm.DB) error) error
	Get(id uuid.UUID) (*Network, error)
	GetByName(orgName string, network string) (*Network, error)
	GetByOrg(orgID uuid.UUID) ([]Network, error)
	// GetByOrgName(orgName string) ([]Network, error)
	// Update(orgId uint, network *Network) error
	Delete(orgName string, network string) error
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

func (n netRepo) GetByName(orgName string, network string) (*Network, error) {
	var ntwk Network

	result := n.Db.GetGormDb().Joins("JOIN orgs on orgs.id=networks.org_id").
		Where("orgs.name=? and networks.name=? and orgs.deleted_at is null",
			orgName, network).Find(&ntwk)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &ntwk, nil
}

func (n netRepo) GetByOrg(orgID uuid.UUID) ([]Network, error) {
	db := n.Db.GetGormDb()
	var networks []Network

	result := db.Where(&Network{OrgId: orgID}).Find(&networks)
	if result.Error != nil {
		return nil, result.Error
	}

	return networks, nil
}

// func (n netRepo) GetByOrgName(orgID uint) ([]Network, error) {

//This gives the result in a single sql query, but fail to distingush between
//	when org does not exist vs when org has no networks, can improve later.
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
func (n netRepo) GetNetworkCount() (int64, error) {
	var count int64
	result := n.Db.GetGormDb().Model(&Network{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (n netRepo) Delete(orgName string, network string) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		txGorm := sql.NewDbFromGorm(tx, pkg.IsDebugMode)
		txr := NewNetRepo(txGorm)

		net, err := txr.GetByName(orgName, network)
		if err != nil {
			return err
		}

		err = tx.Delete(net).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
