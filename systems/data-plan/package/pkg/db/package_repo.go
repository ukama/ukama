package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type PackageRepo interface {
	Add(_package *Package) error
	Get(id uint64) (*Package, error)
	Delete(id uint64) error
	GetByOrg(orgId uint64) ([]Package, error)
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
	result := r.Db.GetGormDb().Create(_package)

	return result.Error
}

func (p *packageRepo) Get(id uint64) (*Package, error) {
	var _package Package

	result := p.Db.GetGormDb().Where("id = ?", id).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

func (p *packageRepo) GetByOrg(orgId uint64) ([]Package, error) {
	var packages []Package
	result := p.Db.GetGormDb().Where(&Package{Org_id: uint(orgId)}).Find(&packages)

	if result.Error != nil {
		return nil, result.Error
	}
	return packages, nil
}

func (r *packageRepo) Delete(packageId uint64) error {
	result := r.Db.GetGormDb().Where("id = ?", packageId).Delete(&Package{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (b *packageRepo) Update(Id uint64, pkg Package) (*Package, error) {
	result := b.Db.GetGormDb().Where(Id).UpdateColumns(pkg)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pkg, nil
}
