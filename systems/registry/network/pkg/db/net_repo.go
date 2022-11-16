package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/registry/network/pkg"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"
)

type NetRepo interface {
	Get(orgName string, network string) (*Network, error)
	GetByOrg(orgID uint) ([]Network, error)
	Add(orgId uint, network string) (*Network, error)
	// Update(orgId uint, network *Network) error
	Delete(orgName string, network string) error
}

type netRepo struct {
	Db sql.Db
}

func NewNetRepo(db sql.Db) NetRepo {
	return &netRepo{
		Db: db,
	}
}

func (n netRepo) Get(orgName string, network string) (*Network, error) {
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

func (n netRepo) GetByOrg(orgID uint) ([]Network, error) {
	db := n.Db.GetGormDb()
	var networks []Network

	//This gives the result in a single sql query, but fail to distingush between
	//	when org does not exist vs when org has no networks, can improve later.
	// result := db.Joins("JOIN orgs on orgs.id=networks.org_id").
	// Where("orgs.name=? and orgs.deleted_at is null", orgName).Debug().Find(&networks)

	result := db.Where(&Network{OrgID: orgID}).Find(&networks)
	if result.Error != nil {
		return nil, result.Error
	}

	return networks, nil
}

func (n netRepo) Add(orgId uint, network string) (*Network, error) {
	db := n.Db.GetGormDb()

	if !validation.IsValidDnsLabelName(network) {
		return nil, fmt.Errorf("invalid name. must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	netw := &Network{
		OrgID: orgId,
		Name:  network,
	}

	db = db.Create(netw)

	return netw, db.Error
}

func (n netRepo) Delete(orgName string, network string) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		txGorm := sql.NewDbFromGorm(tx, pkg.IsDebugMode)
		txr := NewNetRepo(txGorm)

		net, err := txr.Get(orgName, network)
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
