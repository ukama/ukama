package db

import (
	"github.com/ukama/openIoR/services/common/sql"
	"github.com/ukama/openIoR/services/common/ukama"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	AddOrUpdateNode(node *Node) error
	GetNode(nodeId ukama.NodeID) (*Node, error)
	DeleteNode(nodeId ukama.NodeID) error
	ListNodes() (*[]Node, error)
	GetNodeStatus(nodeId ukama.NodeID) (*string, error)
	UpdateNodeProdStatus(node *Node) error
	UpdateNodeStatus(nodeId ukama.NodeID, status string) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) *nodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) AddOrUpdateNode(node *Node) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "part_number", "skew", "mac", "sw_version", "p_sw_version", "assembly_date", "oem_name", "prod_test_status", "prod_report", "status"}),
	}).Create(node)
	return d.Error
}

func (r *nodeRepo) GetNode(nodeId ukama.NodeID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

/* Delete Node  */
func (r *nodeRepo) DeleteNode(nodeId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("node_id = ?", nodeId).Delete(&Node{})
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

func (r *nodeRepo) GetNodeStatus(nodeId ukama.NodeID) (*string, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node.Status, nil
}

/* Update Production status */
func (r *nodeRepo) UpdateNodeProdStatus(node *Node) error {
	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"prod_test_status", "prod_report", "status"}),
	}).Create(node)
	return d.Error
}

/* Update Node status */
func (r *nodeRepo) UpdateNodeStatus(nodeId ukama.NodeID, status string) error {
	node := Node{
		NodeID: nodeId,
		Status: status,
	}

	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(node)
	return d.Error
}
