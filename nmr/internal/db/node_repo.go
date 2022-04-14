package db

import (
	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/sql"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	AddOrUpdateNode(node *Node) error
	GetNode(nodeId string) (*Node, error)
	DeleteNode(nodeId string) error
	ListNodes() (*[]Node, error)
	GetNodeStatus(nodeId string) (*string, error)
	GetNodeMfgStatus(nodeId string) (*string, *[]byte, error)
	UpdateNodeMfgStatus(node *Node) error
	UpdateNodeStatus(nodeId string, status string) error
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
		UpdateAll: true,
	}).Create(node)
	return d.Error
}

func (r *nodeRepo) GetNode(nodeId string) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().First(&node, "node_id = ?", nodeId)
	if result.Error != nil {
		return nil, result.Error
	}
	logrus.Tracef("Read node info for %s with %v. result %+v", node.NodeID, node, result)
	return &node, nil
}

/* Delete Node  */
func (r *nodeRepo) DeleteNode(nodeId string) error {
	result := r.Db.GetGormDb().Unscoped().Where("node_id = ?", nodeId).Delete(&Node{})
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

func (r *nodeRepo) GetNodeStatus(nodeId string) (*string, error) {
	var node Node
	result := r.Db.GetGormDb().Select("status").First(&node, "node_id = ?", nodeId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &node.Status, nil
}

func (r *nodeRepo) GetNodeMfgStatus(nodeId string) (*string, *[]byte, error) {
	var node Node
	result := r.Db.GetGormDb().Select("status", "mfg_report").First(&node, "node_id = ?", nodeId)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return &node.Status, node.MfgReport, nil
}

/* Update Production status */
func (r *nodeRepo) UpdateNodeMfgStatus(node *Node) error {

	result := r.Db.GetGormDb().Model(&Node{}).Where("node_id = ?", node.NodeID).UpdateColumns(node)
	if result.Error != nil {
		return result.Error
	}
	logrus.Tracef("Updated node mfg status for %s with %v. result %+v", node.NodeID, node, result)
	return nil

	// d := r.Db.GetGormDb().Clauses(clause.OnConflict{
	// 	Columns:   []clause.Column{{Name: "node_id"}},
	// 	DoUpdates: clause.AssignmentColumns([]string{"mfg_test_status", "mfg_report", "status"}),
	// }).Create(node)
	// return d.Error
}

/* Update Node status */
func (r *nodeRepo) UpdateNodeStatus(nodeId string, status string) error {
	result := r.Db.GetGormDb().Model(&Node{}).Where("node_id = ?", nodeId).UpdateColumn("status", status)
	if result.Error != nil {
		return result.Error
	}
	logrus.Tracef("Updated node status for %s with %s. result %+v", nodeId, status, result)
	return nil
}
