package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
)

type StateRepo interface {
	Create(state *State, nestedFunc func(*State, *gorm.DB) error) error
	GetByNodeId(nodeId ukama.NodeID) (*State, error)
	Delete(nodeId ukama.NodeID) error
	GetStateHistory(nodeId ukama.NodeID) ([]State, error)
}

type stateRepo struct {
	Db sql.Db
}

func NewStateRepo(db sql.Db) StateRepo {
	return &stateRepo{
		Db: db,
	}
}

func (r *stateRepo) Create(state *State, nestedFunc func(state *State, tx *gorm.DB) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(state, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(state)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (r *stateRepo) GetByNodeId(nodeId ukama.NodeID) (*State, error) {
	var state State
	err := r.Db.GetGormDb().
		Where("node_id = ?", nodeId.String()).
		Order("created_at desc").
		First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *stateRepo) Delete(nodeId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("node_id = ?", nodeId.String()).Delete(&State{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *stateRepo) GetStateHistory(nodeId ukama.NodeID) ([]State, error) {
	var states []State
	err := r.Db.GetGormDb().
		Where("node_id = ?", nodeId.String()).
		Order("created_at desc").
		Find(&states).Error
	if err != nil {
		return nil, err
	}
	return states, nil
}
