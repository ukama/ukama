// This is an example of a repositoryasrRepo
package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

// declare interface so that we can mock it
type PolicyRepo interface {
	Add(policy *Policy) error
	Get(id uuid.UUID) (*Policy, error)
	Delete(id uuid.UUID) error
	Update(id uint, policy *Policy) error
	GetByAsrId(id uint) (*Policy, error)
}

type policyRepo struct {
	db sql.Db
}

func NewPolicyRepo(db sql.Db) *policyRepo {
	return &policyRepo{
		db: db,
	}
}

func (p *policyRepo) Add(policy *Policy) error {

	r := p.db.GetGormDb().Create(policy)
	if r.Error != nil {
		log.Errorf("error creating policy %+v. Error: %v", policy, r.Error)
		return r.Error
	}

	return nil
}

func (p *policyRepo) Get(id *uuid.UUID) (*Policy, error) {
	var policy Policy
	result := p.db.GetGormDb().Where("id = ? AND delete_at = null", id).First(&policy)
	if result.Error != nil {
		log.Errorf("error reading policy %s. Error: %v", id.String(), err)
		return nil, result.Error
	}

	return &policy, nil
}

func (p *policyRepo) GetByAsrId(id uint) (*Policy, error) {
	var policy Policy
	result := p.db.GetGormDb().Where("asr_id = ? AND delete_at = null", id).First(&policy)
	if result.Error != nil {
		log.Errorf("error reading policy for ASR ID %d. Error: %v", id, err)
		return nil, result.Error
	}

	return &policy, nil
}

func (r *policyRepo) Delete(id uuid.UUID, nestedFunc ...func(*gorm.DB) error) error {
	return r.db.ExecuteInTransaction2(func(tx *gorm.DB) *gorm.DB {
		return tx.Where(&Policy{Id: id}).Delete(&Policy{})
	}, nestedFunc...)
}

func (r *policyRepo) Update(id uuid.UUID, newPolicy Policy) error {

	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {

		policy := &Policy{}

		res := tx.Model(&Policy{}).Where("asr_id = ?", id).Find(policy)
		if res.Error != nil {
			if !sql.IsNotFoundError(res.Error) {
				log.Errorf("Error looking for policy for id %d.Error %+v", id, res.Error)
				return res.Error
			}
		}

		if err := tx.Delete(policy).Error; err != nil {
			log.Errorf("Failed to delete policy %+v  for ASR id %d .Error %s", id, policy, err.Error())
			return nil
		}

		if err := tx.Create(newPolicy).Error; err != nil {
			log.Errorf("Failed to create policy %+v for ASR id %d .Error %s", newPolicy, id, err.Error())
			return nil
		}

		return nil
	})

	return err
}
