package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SoftwareManagerRepo interface {
	CreateSoftwareUpdate(Software *Software, nestedFunc func(string, string) error) error
	GetLatestSoftwareUpdate() (*Software, error)
}
type softwareManagerRepo struct {
	Db sql.Db
}

func NewSoftwareManagerRepo(db sql.Db) SoftwareManagerRepo {
	return &softwareManagerRepo{
		Db: db,
	}
}
func (r *softwareManagerRepo) CreateSoftwareUpdate(Software *Software, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc("", "")
			if nestErr != nil {
				return nestErr
			}
		}
		if err := tx.Create(Software).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *softwareManagerRepo) GetLatestSoftwareUpdate() (*Software, error) {
	var Software Software
	err := r.Db.GetGormDb().Order("release_date desc").First(&Software).Error
	return &Software, err
}
