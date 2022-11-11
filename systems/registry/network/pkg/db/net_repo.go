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
	Add(orgId uint, network string) (*Network, error)

	// List gets list of org, networks and node count by type
	// returns map[org][network][nodeType]=nodeCount
	List() (map[string]map[string]map[NodeType]int, error)
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
	db := n.Db.GetGormDb()

	rows, err := db.Raw(`SELECT n.name,n.id, n.org_id, o.name  from networks n
inner join orgs o on n.org_id = o.id
where n.deleted_at IS NULL and o.deleted_at IS NULL and
 n.name = ? and o.name = ?`, network, orgName).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	nt := Network{
		Org: &Org{},
	}

	exist := false

	for rows.Next() {
		exist = true

		err = rows.Scan(&nt.Name, &nt.ID, &nt.OrgID, &nt.Org.Name)

		if err != nil {
			return nil, err
		}
	}

	if !exist {
		return nil, gorm.ErrRecordNotFound
	}

	return &nt, nil
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

func (n netRepo) List() (map[string]map[string]map[NodeType]int, error) {
	db := n.Db.GetGormDb()

	rows, err := db.Raw(`select  o."name" org ,n."name" network ,nd.type , count(n.id) nodes
from orgs o
         inner join networks n  on n.org_id  = o.id
         inner join nodes nd on nd.network_id  = n.id and nd.deleted_at is null
group by o.id , n.id, nd.type`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]map[string]map[NodeType]int)

	for rows.Next() {
		var org, network string
		var nodes int
		var nodeType NodeType

		err = rows.Scan(&org, &network, &nodeType, &nodes)
		if err != nil {
			return nil, err
		}

		if _, ok := result[org]; !ok {
			result[org] = make(map[string]map[NodeType]int)
		}

		if _, ok := result[org][network]; !ok {
			result[org][network] = make(map[NodeType]int)
		}

		result[org][network][nodeType] = nodes
	}

	return result, nil
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

		err = tx.Where("network_id = ?", net.ID).Delete(&Node{}).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
