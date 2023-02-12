package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type PackageRepo interface {
	Add(_package *Package) error
	Get(uuid uuid.UUID) (*Package, error)
	Delete(uuid uuid.UUID) error
	GetByOrg(orgId uuid.UUID) ([]Package, error)
	Update(uuid uuid.UUID, pkg Package) (*Package, error)
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

func (p *packageRepo) Get(uuid uuid.UUID) (*Package, error) {
	var _package Package

	result := p.Db.GetGormDb().Where("uuid = ?", uuid).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

func (p *packageRepo) GetByOrg(orgId uuid.UUID) ([]Package, error) {
	var packages []Package
	result := p.Db.GetGormDb().Where(&Package{Org_id: orgId}).Find(&packages)

	if result.Error != nil {
		return nil, result.Error
	}
	return packages, nil
}

func (r *packageRepo) Delete(uuid uuid.UUID) error {
	result := r.Db.GetGormDb().Where("uuid = ?", uuid).Delete(&Package{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (b *packageRepo) Update(uuid uuid.UUID, pkg Package) (*Package, error) {
	result := b.Db.GetGormDb().Where("uuid = ?", uuid).UpdateColumns(pkg)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pkg, nil
}
