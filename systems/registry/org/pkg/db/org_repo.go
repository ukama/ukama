package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
)

type OrgRepo interface {
	/* Orgs */
	Add(org *Org) error
	Get(id uint) (*Org, error)
	GetByName(name string) (*Org, error)
	GetByOwner(uuid uuid.UUID) ([]Org, error)
	// Update(id uint) error
	// Deactivate(id uint) error
	// Delete(id uint) error

	/* Members */
	AddMember(member *OrgUser) error
	GetMember(orgID uint, userUUID uuid.UUID) (*OrgUser, error)
	GetMembers(orgID uint) ([]OrgUser, error)
	UpdateMember(orgID uint, member *OrgUser) error
	RemoveMember(orgID uint, userUUID uuid.UUID) error
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org) (err error) {
	if !validation.IsValidDnsLabelName(org.Name) {
		return fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	d := r.Db.GetGormDb().Create(org)

	return d.Error
}

func (r *orgRepo) Get(id uint) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByName(name string) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().Where(&Org{Name: name}).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByOwner(uuid uuid.UUID) ([]Org, error) {
	var orgs []Org

	result := r.Db.GetGormDb().Where(&Org{Owner: uuid}).Find(&orgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgs, nil
}

func (r *orgRepo) AddMember(member *OrgUser) error {
	d := r.Db.GetGormDb().Create(member)

	return d.Error
}

func (r *orgRepo) GetMember(orgID uint, userUUID uuid.UUID) (*OrgUser, error) {
	var member OrgUser

	result := r.Db.GetGormDb().Where("org_id = ? And uuid = ?", orgID, userUUID).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}

	return &member, nil
}

func (r *orgRepo) GetMembers(orgID uint) ([]OrgUser, error) {
	var members []OrgUser

	result := r.Db.GetGormDb().Where(&OrgUser{OrgID: uint(orgID)}).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}

	return members, nil
}

func (r *orgRepo) UpdateMember(orgID uint, member *OrgUser) error {
	d := r.Db.GetGormDb().Clauses(clause.Returning{}).Where("org_id = ? And uuid = ?", member.OrgID, member.Uuid).Updates(member)
	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return d.Error
}

func (r *orgRepo) RemoveMember(orgID uint, userUUID uuid.UUID) error {
	var member OrgUser

	// d := r.Db.GetGormDb().Clauses(clause.Returning{}).Where("org_id = ? And uuid = ?", orgID, userUUID).Delete(&member)
	d := r.Db.GetGormDb().Where("org_id = ? And uuid = ?", orgID, userUUID).Delete(&member)
	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return d.Error
}
