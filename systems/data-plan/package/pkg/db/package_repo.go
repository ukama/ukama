package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PackageRepo interface {
	Add(_package *Package) error
	Get(uuid uuid.UUID) (*Package, error)
	GetDetails(uuid.UUID) (*Package, error)
	Delete(uuid uuid.UUID) error
	GetByOrg(orgId uuid.UUID) ([]Package, error)
	Update(uuid uuid.UUID, pkg *Package) error
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

	result := p.Db.GetGormDb().Preload("PackageRate").Where("uuid = ?", uuid).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

func (p *packageRepo) GetDetails(uuid uuid.UUID) (*Package, error) {
	var _package Package

	result := p.Db.GetGormDb().Preload(clause.Associations).Where("uuid = ?", uuid).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

func (p *packageRepo) GetByOrg(orgId uuid.UUID) ([]Package, error) {
	var packages []Package
	result := p.Db.GetGormDb().Preload("PackageRate").Where(&Package{OrgId: orgId}).Find(&packages)

	if result.Error != nil {
		return nil, result.Error
	}
	return packages, nil
}

func (r *packageRepo) Delete(uuid uuid.UUID) error {
	p := &Package{}
	result := r.Db.GetGormDb().Model(&Package{}).Where("uuid=?", uuid).Delete(p)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (b *packageRepo) Update(uuid uuid.UUID, pkg *Package) error {
	d := b.Db.GetGormDb().Clauses(clause.Returning{}).Where("uuid = ?", uuid).Updates(pkg)
	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	//https://stackoverflow.com/questions/65683156/updates-doesnt-seem-to-update-the-associations

	return d.Error

}
