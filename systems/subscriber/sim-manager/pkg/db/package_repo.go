package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type PackageRepo interface {
	Add(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error
	Delete(packageID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) PackageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (u *packageRepo) Add(pkg *Package, nestedFunc func(pkg *Package, tx *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(pkg)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(pkg, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (u *packageRepo) Delete(packageID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Package{}, packageID)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(packageID, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
