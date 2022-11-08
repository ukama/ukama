package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/validation"

	"github.com/ukama/ukama/systems/common/sql"
)

type OrgRepo interface {
	Add(org *Org) error
	Get(id int) (*Org, error)
	GetByOwner(uuid uuid.UUID) ([]Org, error)
	// Update(id int) error
	// Deactivate(id int) error
	// Delete(id int) error

	AddMember(org *Org, user *User) (*OrgUser, error)
	// GetMember()
	// GetMembers()

	// DeactivateMember()
	// RemoveMember()
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	_ = db.GetGormDb().SetupJoinTable(&Org{}, "Members", &OrgUser{})
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

func (r *orgRepo) Get(id int) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().First(&org, id)
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

func (r *orgRepo) AddMember(org *Org, user *User) (*OrgUser, error) {
	if !validation.IsValidDnsLabelName(org.Name) {
		return nil, fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	member := &OrgUser{
		OrgID:  org.ID,
		UserID: user.ID,
		Uuid:   user.Uuid,
	}

	d := r.Db.GetGormDb().Create(member)

	return member, d.Error
}

// func (r *orgRepo) Delete(name string) error {
// re,turn r.Db.GetGormDb().Delete(&Org{}, "name = ?", name).Error
// }
