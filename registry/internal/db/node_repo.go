package db

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukamaX/common/sql"
	"github.com/ukama/ukamaX/common/ukama"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	Add(node *Node) error
	Get(id ukama.NodeID) (*Node, error)
	GetByOrg(orgName string, ownerId uuid.UUID) ([]Node, error)
	Delete(id ukama.NodeID) error
	Update(id ukama.NodeID, state NodeState) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) Add(node *Node) error {
	d := r.Db.GetGormDb().Create(node)
	return d.Error
}

func (r *nodeRepo) Get(id ukama.NodeID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id=?", id.String())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

func (r *nodeRepo) GetByOrg(orgName string, ownerId uuid.UUID) ([]Node, error) {
	var nodes []Node

	err := r.Db.GetGormDb().Joins("inner join orgs ON orgs.id = nodes.org_id").Where("orgs.name=? and orgs.owner=?", orgName, ownerId).Find(&nodes)
	if err.Error != nil {
		return nil, err.Error
	}
	return nodes, nil
}

func (r *nodeRepo) Delete(id ukama.NodeID) error {
	res := r.Db.GetGormDb().Delete(&Node{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *nodeRepo) Update(id ukama.NodeID, state NodeState) error {

	result := r.Db.GetGormDb().Where("node_id=?", id.String()).Updates(Node{State: state})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
