package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type SoftwareManagerRepo interface {
	CreateSoftware(Software *Software, nestedFunc func(string, string) error) error
	ReadSoftware(id uuid.UUID) (*Software, error)
	ListSoftwares() ([]*Software, error)
	GetLatestSoftware() (*Software, error)
}
type softwareManagerRepo struct {
	Db sql.Db
}

func NewSoftwareManagerRepo(db sql.Db) SoftwareManagerRepo {
	return &softwareManagerRepo{
		Db: db,
	}
}
func (r *softwareManagerRepo) CreateSoftware(Software *Software, nestedFunc func(string, string) error) error {
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
func (r *softwareManagerRepo) ReadSoftware(id uuid.UUID) (*Software, error) {
	var Software Software
	err := r.Db.GetGormDb().Where("id = ?", id).First(&Software).Error
	return &Software, err
}
func (r *softwareManagerRepo) ListSoftwares() ([]*Software, error) {
	var Softwares []*Software
	err := r.Db.GetGormDb().Find(&Softwares).Error
	return Softwares, err
}
func (r *softwareManagerRepo) GetLatestSoftware() (*Software, error) {
	var Software Software
	err := r.Db.GetGormDb().Order("release_date desc").First(&Software).Error
	return &Software, err
}
