package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PackageRepo interface {
	Add(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error
	Get(packageID uuid.UUID) (*Package, error)
	GetBySim(simID uuid.UUID) ([]Package, error)
	GetOverlap(*Package) ([]Package, error)
	Update(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error
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

func (p *packageRepo) Add(pkg *Package, nestedFunc func(pkg *Package, tx *gorm.DB) error) error {
	err := p.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
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

func (p *packageRepo) Get(packageID uuid.UUID) (*Package, error) {
	var pkg Package

	result := p.Db.GetGormDb().First(&pkg, packageID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pkg, nil
}

func (p *packageRepo) GetBySim(simID uuid.UUID) ([]Package, error) {
	var packages []Package

	result := p.Db.GetGormDb().Where(&Package{SimID: simID}).Find(&packages)
	if result.Error != nil {
		return nil, result.Error
	}

	return packages, nil
}

func (p *packageRepo) GetOverlap(pkg *Package) ([]Package, error) {
	var packages []Package

	result := p.Db.GetGormDb().Where(&Package{SimID: pkg.SimID}).Find(&packages,
		"end_date >= ? AND start_date <= ?", pkg.StartDate, pkg.EndDate)

	if result.Error != nil {
		return nil, result.Error
	}

	return packages, nil
}

func (p *packageRepo) Update(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error {
	err := p.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(pkg, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Clauses(clause.Returning{}).Updates(pkg)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (p *packageRepo) Delete(packageID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := p.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
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
