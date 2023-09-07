package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	AddOrUpdate(node *Node) error
	Get(nodeId ukama.NodeID) (*Node, error)
	Delete(nodeId ukama.NodeID) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) *nodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) AddOrUpdate(node *Node) error {
	d := r.Db.GetGormDb().Create(node)
	return d.Error
}

func (r *nodeRepo) Get(nodeId ukama.NodeID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

func (r *nodeRepo) Delete(nodeId ukama.NodeID) error {
	var node Node
	result := r.Db.GetGormDb().Unscoped().Preload(clause.Associations).Delete(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("%s node missing", nodeId.String())
}
