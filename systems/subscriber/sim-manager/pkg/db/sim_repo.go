package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SimRepo interface {
	Add(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error
	Get(simID uuid.UUID) (*Sim, error)
	GetBySubscriber(subscriberID uuid.UUID) ([]Sim, error)
	GetByNetwork(networkID uuid.UUID) ([]Sim, error)
	Update(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error
	Delete(simID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) SimRepo {
	return &simRepo{
		Db: db,
	}
}

func (s *simRepo) Add(sim *Sim, nestedFunc func(sim *Sim, tx *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(sim)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(sim, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (s *simRepo) Get(simID uuid.UUID) (*Sim, error) {
	var sim Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Preload("Package", "is_active is true").First(&sim, simID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &sim, nil
}

func (s *simRepo) GetBySubscriber(subscriberID uuid.UUID) ([]Sim, error) {
	var sims []Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{SubscriberID: subscriberID}).Preload("Package", "is_active is true").Find(&sims)
	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

func (s *simRepo) GetByNetwork(networkID uuid.UUID) ([]Sim, error) {
	var sims []Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{NetworkID: networkID}).Preload("Package", "is_active is true").Find(&sims)
	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

// Update sim modified non-empty fields provided by Sim struct
func (s *simRepo) Update(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Returning{}).Updates(sim)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(sim, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (s *simRepo) Delete(simID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Sim{}, simID)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(simID, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
