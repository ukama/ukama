package db

import (
	"github.com/ukama/ukama/services/common/sql"
	"gorm.io/gorm/clause"
)

type VNodeRepo interface {
	Insert(nodeId string, status string) error
	Update(nodeId string, status string) error
	PowerOn(nodeId string) error
	PowerOff(nodeId string) error
	GetInfo(nodeId string) (*VNode, error)
	List() (*[]VNode, error)
	Delete(nodeId string) error
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
func (r *vNodeRepo) Update(nodeId string, status string) error {
	d := r.Db.GetGormDb().Where("node_id = ?", nodeId).Updates(VNode{Status: status})
	return d.Error
}

/* Update is used when we know the node id */
func (r *vNodeRepo) Insert(nodeId string, status string) error {
	vNode := &VNode{
		NodeID: nodeId,
		Status: status,
	}

	d := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoNothing: true,
	}).Create(vNode)
	return d.Error
}

/* PowerOn */
func (r *vNodeRepo) PowerOn(nodeId string) error {
	return r.Update(nodeId, VNodeOn.String())
}

/* PowerOn */
func (r *vNodeRepo) PowerOff(nodeId string) error {
	return r.Update(nodeId, VNodeOff.String())
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

/* Delete VirtNode info */
func (r *vNodeRepo) Delete(nodeId string) error {
	result := r.Db.GetGormDb().Unscoped().Where("node_id = ?", nodeId).Delete(&VNode{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
