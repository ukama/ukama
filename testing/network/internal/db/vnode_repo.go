package db

import (
	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type VNodeRepo interface {
	Upsert(nodeId string, status string) error
	PowerOn(nodeId string) error
	PowerOff(nodeId string) error
	GetInfo(nodeId string) (*VNode, error)
	List() (*[]VNode, error)
}

type vNodeRepo struct {
	Db sql.Db
}

func NewVNodeRepo(db sql.Db) *vNodeRepo {
	return &vNodeRepo{
		Db: db,
	}
}

/* Upsert is used when we know the node id */
func (r *vNodeRepo) Upsert(nodeId string, status string) error {
	vNode := &VNode{
		NodeID: nodeId,
		Status: status,
	}

	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		UpdateAll: true,
	}).Create(vNode)
	return d.Error
}

/* PowerOn */
func (r *vNodeRepo) PowerOn(nodeId string) error {
	return r.Upsert(nodeId, VNodeOn.String())
}

/* PowerOn */
func (r *vNodeRepo) PowerOff(nodeId string) error {
	return r.Upsert(nodeId, VNodeOff.String())
}

/* Get VirtNode info */
func (r *vNodeRepo) GetInfo(nodeId string) (*VNode, error) {
	vNode := VNode{}
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&vNode, "node_id = ?", nodeId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vNode, nil
}

/* List all Modules */
func (r *vNodeRepo) List() (*[]VNode, error) {
	vNodes := []VNode{}

	result := r.Db.GetGormDb().Preload(clause.Associations).Find(&vNodes)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	} else {
		return &vNodes, nil
	}
}
