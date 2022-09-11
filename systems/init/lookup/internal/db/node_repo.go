package db

import (
	"github.com/ukama/ukama/services/common/sql"
	"github.com/ukama/ukama/services/common/ukama"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	AddOrUpdate(node *Node) error
	Get(nodeId ukama.NodeID) (*Node, error)
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
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "lower(node_id::text)", Raw: true}},
		DoUpdates: clause.AssignmentColumns([]string{"org_id"}),
	}).Create(node)
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
