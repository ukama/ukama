package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type PackageRepo interface {
	Add(_package *Package) error
	Get(orgId, id uint64) ([]Package, error)
	Delete(orgId, id uint64) error
	Update(Id uint64, pkg Package) (*Package, error)
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) *packageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (r *packageRepo) Add(_package *Package) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(_package)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (p *packageRepo) Get(orgId, id uint64) ([]Package, error) {
	var packages []Package
	result := p.Db.GetGormDb().Where(&Package{Org_id: uint(orgId), Model: gorm.Model{ID: uint(id)}}).Find(&packages)

	if result.Error != nil {
		return nil, result.Error
	}
	return packages, nil
}

func (r *packageRepo) Delete(orgId uint64, packageId uint64) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND org_id = ?", packageId, orgId).Delete(&Package{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			result.Error = gorm.ErrRecordNotFound
			return result.Error
		}

		return nil
	})

	return err
}

func (b *packageRepo) Update(Id uint64, pkg Package) (*Package, error) {

	result := b.Db.GetGormDb().Where(Id).UpdateColumns(pkg)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pkg, nil
}
