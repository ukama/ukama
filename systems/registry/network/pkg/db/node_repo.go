package db

import (
	sql2 "database/sql"
	"strings"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MaxAttachedNodes = 2

// NodeID must be lowercase
type NodeRepo interface {
	Add(node *Node, nestedFunc ...func() error) error
	Get(id ukama.NodeID) (*Node, error)
	GetByOrg(orgName string) ([]Node, error)
	Delete(id ukama.NodeID, nestedFunc ...func() error) error
	Update(id ukama.NodeID, node *NodeAttributes, nestedFunc ...func() error) error
}

type NodeAttributes struct {
	State *NodeState
	Name  *string
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

	result := r.Db.GetGormDb().Preload("Network.Org").Preload(clause.Associations).First(&node, "node_id=?", id.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

func (r *nodeRepo) GetByOrg(orgName string) ([]Node, error) {
	db := r.Db.GetGormDb()

	rows, err := db.Raw(`select n.* from nodes n
									inner join networks nw on nw.id = n.network_id	
									inner join orgs o ON o.id = nw.org_id
								where o.name=? and n.deleted_at is null and nw.deleted_at is null`, orgName).Rows()
	if err != nil {
		return nil, err
	}

	nodes, err := r.mapNodes(rows, db)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepo) mapNodes(rows *sql2.Rows, db *gorm.DB) ([]Node, error) {
	var nodes []Node
	defer rows.Close()

	for rows.Next() {
		var node Node

		// Node columns are mapped correctly
		err := db.ScanRows(rows, &node)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (r *nodeRepo) Delete(id ukama.NodeID, nestedFunc ...func() error) error {
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		d := tx.Delete(&Node{}, "node_id = ?", id.StringLowercase())
		if d.Error != nil {
			return d
		}

		if d.RowsAffected == 0 {
			d.Error = gorm.ErrRecordNotFound
			return d
		}

		return d
	}, nestedFunc...)

	return err
}

// Update updated node with `id`. Only fields that are not nil are updated
func (r *nodeRepo) Update(id ukama.NodeID, node *NodeAttributes, nestedFunc ...func() error) error {
	var rowsAffected int64

	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		nd := Node{}
		if node.State != nil {
			nd.State = *node.State
		}
		if node.Name != nil {
			nd.Name = *node.Name
		}

		result := tx.Where("node_id=?", id.StringLowercase()).Updates(nd)
		rowsAffected = result.RowsAffected

		return result
	}, nestedFunc...)

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return err
}
