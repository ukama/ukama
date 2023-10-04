package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type InvitationRepo interface {
	Add(invitation *Invitation, nestedFunc func(*Invitation, *gorm.DB) error) error
	Get(id uuid.UUID) (*Invitation, error)
	GetByOrg(org string) ([]*Invitation, error)
	UpdateStatus(id uuid.UUID, status uint8) error
	Delete(id uuid.UUID, nestedFunc func(string, string) error) error
	GetInvitationByEmail(email string) (*Invitation, error)
}

type invitationRepo struct {
	Db sql.Db
}

func NewInvitationRepo(db sql.Db) InvitationRepo {
	return &invitationRepo{
		Db: db,
	}
}

func (r *invitationRepo) GetInvitationByEmail(email string) (*Invitation, error) {
	var invitation Invitation
	err := r.Db.GetGormDb().Where("email = ?", email).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepo) Add(invitation *Invitation, nestedFunc func(*Invitation, *gorm.DB) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(invitation, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		if err := tx.Create(invitation).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *invitationRepo) Get(id uuid.UUID) (*Invitation, error) {
	var invitation Invitation
	err := r.Db.GetGormDb().First(&invitation, id).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepo) GetByOrg(org string) ([]*Invitation, error) {
	var invitations []*Invitation
	err := r.Db.GetGormDb().Where("org = ?", org).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *invitationRepo) Delete(id uuid.UUID, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc("", "")
			if nestErr != nil {
				return nestErr
			}
		}

		if err := tx.Delete(&Invitation{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *invitationRepo) UpdateStatus(id uuid.UUID, status uint8) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Invitation{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
