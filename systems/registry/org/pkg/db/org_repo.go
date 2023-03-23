package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
)

type OrgRepo interface {
	/* Orgs */
	Add(org *Org, nestedFunc func(*Org, *gorm.DB) error) error
	Get(id uuid.UUID) (*Org, error)
	GetByName(name string) (*Org, error)
	GetByOwner(uuid uuid.UUID) ([]Org, error)
	// Update(id uint) error
	// Deactivate(id uint) error
	// Delete(id uint) error

	/* Members */
	AddMember(member *OrgUser) error
	GetMember(orgID uuid.UUID, userUUID uuid.UUID) (*OrgUser, error)
	GetMembers(orgID uuid.UUID) ([]OrgUser, error)
	UpdateMember(orgID uuid.UUID, member *OrgUser) error
	RemoveMember(orgID uuid.UUID, userUUID uuid.UUID) error
	GetOrgCount() (int64, int64, error)
	GetMemberCount(orgID uuid.UUID) (int64, int64, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org, nestedFunc func(*Org, *gorm.DB) error) (err error) {
	if !validation.IsValidDnsLabelName(org.Name) {
		return fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err = r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(org, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		d := tx.Create(org)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return err
}

func (r *orgRepo) Get(id uuid.UUID) (*Org, error) {
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

func (r *orgRepo) GetMember(orgID uuid.UUID, userUUID uuid.UUID) (*OrgUser, error) {
	var member OrgUser

	result := r.Db.GetGormDb().Where("org_id = ? And uuid = ?", orgID, userUUID).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}

	return &member, nil
}

func (r *orgRepo) GetMembers(orgID uuid.UUID) ([]OrgUser, error) {
	var members []OrgUser

	result := r.Db.GetGormDb().Where(&OrgUser{OrgId: orgID}).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}

	return members, nil
}

func (r *orgRepo) UpdateMember(orgID uuid.UUID, member *OrgUser) error {
	d := r.Db.GetGormDb().Clauses(clause.Returning{}).Where("org_id = ? And uuid = ?", member.OrgId, member.Uuid).Updates(member)
	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return d.Error
}

func (r *orgRepo) RemoveMember(orgID uuid.UUID, userUUID uuid.UUID) error {
	var member OrgUser

	// d := r.Db.GetGormDb().Clauses(clause.Returning{}).Where("org_id = ? And uuid = ?", orgID, userUUID).Delete(&member)
	d := r.Db.GetGormDb().Where("org_id = ? And uuid = ?", orgID, userUUID).Delete(&member)
	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return d.Error
}

func (r *orgRepo) GetOrgCount() (int64, int64, error) {
	var activeOrgCount int64
	var deactiveOrgCount int64

	result := r.Db.GetGormDb().Model(&Org{}).Where("deactivated = ?", false).Count(&activeOrgCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = r.Db.GetGormDb().Model(&Org{}).Where("deactivated = ?", true).Count(&deactiveOrgCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	return activeOrgCount, deactiveOrgCount, nil
}

func (r *orgRepo) GetMemberCount(orgID uuid.UUID) (int64, int64, error) {
	var activeMemberCount int64
	var deactiveMemberCount int64

	result := r.Db.GetGormDb().Model(&OrgUser{}).Where("org_id = ? AND deactivated = ?", orgID, false).Count(&activeMemberCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = r.Db.GetGormDb().Model(&OrgUser{}).Where("org_id = ? AND deactivated = ?", orgID, true).Count(&deactiveMemberCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	return activeMemberCount, deactiveMemberCount, nil
}
