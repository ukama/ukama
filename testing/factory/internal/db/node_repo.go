package db

import (
	ukama "github.com/ukama/ukama/services/common/ukama"

	"github.com/ukama/ukama/services/common/sql"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	Add(node *Node) error
	UpdateNodeStatus(id string, status string) error
	GetNodeByNodeId(nodeId ukama.NodeID) (*Node, error)
	ListNodes() (*[]Node, error)
	DeleteNodeByNodeId(nodeId ukama.NodeID) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) *nodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

/* Add Virtual Node */
func (r *nodeRepo) Add(node *Node) error {

	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "name", "type", "status", "created_at", "updated_at"}),
	}).Create(node)
	return d.Error
}

/* Update Node Status */
func (r *nodeRepo) UpdateNodeStatus(id string, status string) error {
	d := r.Db.GetGormDb().Where("id = ?", id).Updates(Node{Status: status})
	return d.Error
}

/* Read nodes by Node Id */
func (r *nodeRepo) GetNodeByNodeId(nodeId ukama.NodeID) (*Node, error) {
	var node Node

	result := r.Db.GetGormDb().Where("id = ?", nodeId).Find(&node)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Debugf("Result is %+v", result)
	if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return &node, nil
	}
}

/* Delete Node  */
func (r *nodeRepo) DeleteNodeByNodeId(nodeId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("id = ?", nodeId).Delete(&Node{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* List all nodes */
func (r *nodeRepo) ListNodes() (*[]Node, error) {
	var nodes []Node

	result := r.Db.GetGormDb().Find(&nodes)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return &nodes, nil
	}
}
