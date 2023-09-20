package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type HealthRepo interface {
	StoreRunningAppsInfo(health *Health, nestedFunc func(string, string) error) error
	GetRunningAppsInfo() ([]*Health, error)
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

func (r *healthRepo) GetRunningAppsInfo() ([]*Health, error) {
	var healths []*Health
	err := r.Db.GetGormDb().Find(&healths).Error
	return healths, err
}
