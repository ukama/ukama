package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NodeStateRepo interface {
	Add(nodeState *NodeState) error
	GetNodeState(nodeId string) (*NodeState, error)
	GetAllNodeStates() ([]*NodeState, error)
}

type nodeStateRepo struct {
	Db sql.Db
}

func NewNodeStateRepo(db sql.Db) NodeStateRepo {
	return &nodeStateRepo{
		Db: db,
	}
}

func (r *nodeStateRepo) Add(nodeState *NodeState) error {
	d := r.Db.GetGormDb().Create(nodeState)
	return d.Error
}

func (r *nodeStateRepo) GetNodeState(nodeId string) (*NodeState, error) {
	var nodeState *NodeState

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if nodeId != "" {
		tx = tx.Where("node_id = ?", nodeId)
	} else {
		return nil, fmt.Errorf("invalid node id %s", nodeId)
	}

	result := tx.First(&nodeState)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return nodeState, nil
}

func (r *nodeStateRepo) GetAllNodeStates() ([]*NodeState, error) {
	var nodeStates []*NodeState

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	result := tx.Find(&nodeStates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return nodeStates, nil
}

