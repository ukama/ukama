package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm"
)

type HealthRepo interface {
	StoreRunningAppsInfo(health *Health, nestedFunc func(string, string) error) error
	GetRunningAppsInfo(nodeId ukama.NodeID) (*Health, error)
}
type healthRepo struct {
	Db sql.Db
}

func NewHealthRepo(db sql.Db) HealthRepo {
	return &healthRepo{
		Db: db,
	}
}
func (r *healthRepo) StoreRunningAppsInfo(health *Health, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc("", "")
			if nestErr != nil {
				return nestErr
			}
		}
		if err := tx.Create(health).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}



func (r *healthRepo) GetRunningAppsInfo(nodeId ukama.NodeID) (*Health, error) {
	var healths Health
	result := r.Db.GetGormDb().Where("node_id = ?", nodeId).
		Preload("System").
		Preload("Capps.Resources").
		First(&healths)
	if result.Error != nil {
		return nil, result.Error
	}

	return &healths, nil
}
