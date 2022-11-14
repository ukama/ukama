package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type PackageRepo interface {
	GetPackage(packageId ,OrgId uint64 ) (*Package, error)
	GetPackages(orgId uint64) ([]Package, error)
	AddPackage(_package *Package, nestedFunc ...func() error) error
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

func (r *packageRepo) AddPackage(_package *Package, nestedFunc ...func() error) error {
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(_package)
	}, nestedFunc...)

	return err
}
func (p *packageRepo) GetPackage(packageId,OrgId uint64) (*Package, error) {
	_package := &Package{}
	result :=p.Db.GetGormDb().Where("Id = ? AND Org_id = ?",packageId , OrgId).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}
	return _package, nil
}

func (p *packageRepo) GetPackages( orgId uint64) ([]Package, error) {
	var packages []Package
	result :=p.Db.GetGormDb().Where("Org_id = ?", orgId).Find(&packages)

	if result.Error != nil {
		return nil, result.Error
	}

	return packages, nil
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
