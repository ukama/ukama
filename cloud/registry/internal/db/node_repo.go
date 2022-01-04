package db

import (
	sql2 "database/sql"
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukamaX/common/sql"
	"github.com/ukama/ukamaX/common/ukama"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// NodeID must be lowercase
type NodeRepo interface {
	Add(node *Node, nestedFunc ...func() error) error
	Get(id ukama.NodeID) (*Node, error)
	GetByOrg(orgName string) ([]Node, error)
	Delete(id ukama.NodeID, nestedFunc ...func() error) error
	Update(id ukama.NodeID, state NodeState, nestedFunc ...func() error) error
	GetByUser(ownerId uuid.UUID) ([]Node, error)
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) Add(node *Node, nestedFunc ...func() error) error {
	node.NodeID = strings.ToLower(node.NodeID)
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(node)
	}, nestedFunc...)

	return err
}

func (r *nodeRepo) Get(id ukama.NodeID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id=?", id.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

func (r *nodeRepo) GetByOrg(orgName string) ([]Node, error) {

	db := r.Db.GetGormDb()
	rows, err := db.Raw(`select * from nodes 
									inner join orgs ON orgs.id = nodes.org_id
									where orgs.name=? and nodes.deleted_at is null`, orgName).Rows()
	if err != nil {
		return nil, err
	}
	nodes, err := r.mapNodesToOrgs(rows, db)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepo) GetByUser(ownerId uuid.UUID) ([]Node, error) {
	db := r.Db.GetGormDb()
	rows, err := db.Raw(`select * from nodes 
									inner join orgs ON orgs.id = nodes.org_id
									where orgs.owner=?`, ownerId).Rows()
	if err != nil {
		return nil, err
	}
	nodes, err := r.mapNodesToOrgs(rows, db)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepo) mapNodesToOrgs(rows *sql2.Rows, db *gorm.DB) ([]Node, error) {
	var nodes []Node
	defer rows.Close()

	for rows.Next() {
		var node Node
		var org Org
		err := db.ScanRows(rows, &node)
		if err != nil {
			return nil, err
		}
		err = db.ScanRows(rows, &org)
		if err != nil {
			return nil, err
		}
		node.Org = &org
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (r *nodeRepo) Delete(id ukama.NodeID, nestedFunc ...func() error) error {
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(&Node{}, "node_id = ?", id.StringLowercase())
	}, nestedFunc...)

	return err
}

func (r *nodeRepo) Update(id ukama.NodeID, state NodeState, nestedFunc ...func() error) error {
	var rowsAffected int64
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		result := tx.Where("node_id=?", id.StringLowercase()).Updates(Node{State: state})
		rowsAffected = result.RowsAffected
		return result
	}, nestedFunc...)

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}
