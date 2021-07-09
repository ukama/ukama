package db

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	Add(node *Node) error
	Get(uuid uuid.UUID) (*Node, error)
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

func (r *nodeRepo) Get(uuid uuid.UUID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, uuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}
