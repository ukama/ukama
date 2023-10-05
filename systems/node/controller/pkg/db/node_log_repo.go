package db

import (
	"errors"

	"github.com/ukama/ukama/systems/common/sql"
)

type NodeLogRepo interface {
	Get(nodeId string) (*NodeLog, error)
	Add(nodeLog string) error
}


type nodeLogRepo struct {
	Db sql.Db
}

func NewNodeLogRepo(db sql.Db) NodeLogRepo {
	return &nodeLogRepo{
		Db: db,
	}	
}

func (r *nodeLogRepo) Add(nodeId string) error {
	var nodeLog NodeLog
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).First(&nodeLog).Error; err != nil {
		if err := r.Db.GetGormDb().Create(&NodeLog{NodeId: nodeId}).Error; err != nil {
			return err 
		}
	} else {
		return errors.New("Duplicate record: a record with the same nodeId already exists")
	}
	return nil
}

func (r *nodeLogRepo) Get(nodeId string) (*NodeLog, error) {
	var nodeLog NodeLog
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).First(&nodeLog).Error; err != nil {
		return nil, err
	}
	return &nodeLog, nil
}