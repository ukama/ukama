package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type PackageRepo interface {
	Add(_package *Package, nestedFunc ...func() error) error
	Get(orgId, id uint64) ([]Package, error)
	Delete(orgId, id uint64) (*Package, error)
	Update(id uint64, pkg Package) (*Package, error)
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) *packageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (r *packageRepo) Add(_package *Package, nestedFunc ...func() error) error {
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(_package)
	}, nestedFunc...)

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

func (p *packageRepo) Delete(orgId, packageId uint64) (*Package, error) {
	_package := &Package{}
	result := p.Db.GetGormDb().Where("id = ? AND org_id = ?", packageId, orgId).Delete(_package)
	if result.Error != nil {
		return nil, result.Error
	}

	return _package, nil
}

func (b *packageRepo) Update(Id uint64, pkg Package) (*Package, error) {

	result := b.Db.GetGormDb().Where(Id).UpdateColumns(pkg)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pkg, nil
}
