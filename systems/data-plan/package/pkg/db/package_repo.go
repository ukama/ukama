package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type PackageRepo interface {
	GetPackage(Id uint64) (*Package, error)
	GetPackages() ([]Package, error)
	CreatePackage(Package) (Package, error)
	DeletePackage(Id uint64) (*Package, error)
	UpdatePackage(Id uint64, pkg Package) (*Package, error)
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) *packageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (p *packageRepo) GetPackage(packageId uint64) (*Package, error) {
	_package := &Package{}
	result := p.Db.GetGormDb().First(_package, "Id=?", packageId)
	if result.Error != nil {
		return nil, result.Error
	}
	return _package, nil
}

func (p *packageRepo) GetPackages() ([]Package, error) {
	var packages []Package
	result := p.Db.GetGormDb().Find(&packages)
	if result.Error != nil {
		return nil, result.Error
	}

	return packages, nil
}

func (b *packageRepo) CreatePackage(newPackage Package) (Package, error) {
	_package := Package{}

	result := b.Db.GetGormDb().Create(&newPackage)
	if result.Error != nil {
		return Package{}, result.Error
	}

	return _package, nil
}

func (p *packageRepo) DeletePackage(packageId uint64) (*Package, error) {
	_package := &Package{}
	result := p.Db.GetGormDb().Delete(_package, "id", packageId)
	if result.Error != nil {
		return nil, result.Error
	}
	return _package, nil
}

func (b *packageRepo) UpdatePackage(Id uint64, pkg Package) (*Package, error) {

	result := b.Db.GetGormDb().Where(Id).UpdateColumns(pkg)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pkg, nil
}
