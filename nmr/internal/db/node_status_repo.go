package db

import (
	"github.com/ukama/openIoR/services/common/sql"
	"github.com/ukama/openIoR/services/common/ukama"
	"gorm.io/gorm/clause"
)

type NodeStatusRepo interface {
	AddNodeStatus(node *NodeStatus) error
	GetNodeStatus(nodeId ukama.NodeID) (*NodeStatus, error)
	DeleteNodeStatus(nodeId ukama.NodeID) error
	ListNodeStatus() (*[]NodeStatus, error)
	UpdateNodeProdStatus(node *NodeStatus) error
	UpdateNodeStatus(node *NodeStatus) error
}

type nodeStatusRepo struct {
	Db sql.Db
}

func NewNodeStatusRepo(db sql.Db) *nodeStatusRepo {
	return &nodeStatusRepo{
		Db: db,
	}
}

func (r *nodeStatusRepo) AddNodeStatus(node *NodeStatus) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report", "status"}),
	}).Create(node)
	return d.Error
}

func (r *nodeStatusRepo) GetNodeStatus(nodeId ukama.NodeID) (*NodeStatus, error) {
	var node NodeStatus
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

/* Delete Node  */
func (r *nodeStatusRepo) DeleteNodeStatus(nodeId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("node_id = ?", nodeId).Delete(&NodeStatus{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/* List all nodes */
func (r *nodeStatusRepo) ListNodeStatus() (*[]NodeStatus, error) {
	var nodes []NodeStatus

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

/* Update Production status */
func (r *nodeStatusRepo) UpdateNodeProdStatus(node *NodeStatus) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report", "status"}),
	}).Create(node)
	return d.Error
}

/* Update Node status */
func (r *nodeStatusRepo) UpdateNodeStatus(node *NodeStatus) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(node)
	return d.Error
}
