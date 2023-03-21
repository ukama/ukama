package db

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
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
	GetSimMetrics() (int64, int64, int64, int64, error)
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
		if nestedFunc != nil {
			nestErr := nestedFunc(sim, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		log.Info("Adding sim", sim)
		result := tx.Create(sim)
		if result.Error != nil {
			return result.Error
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

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{SubscriberId: subscriberID}).Preload("Package", "is_active is true").Find(&sims)
	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

func (s *simRepo) GetByNetwork(networkID uuid.UUID) ([]Sim, error) {
	var sims []Sim

	result := s.Db.GetGormDb().Model(&Sim{}).Where(&Sim{NetworkId: networkID}).Preload("Package", "active is true").Find(&sims)
	if result.Error != nil {
		return nil, result.Error
	}

	return sims, nil
}

// Update package modified non-empty fields provided by Package struct
func (s *simRepo) Update(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(sim, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Clauses(clause.Returning{}).Updates(sim)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if result.Error != nil {
			return result.Error
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
func (s *simRepo) GetSimMetrics() (simsCount, activeCount, deactiveCount, terminatedCount int64, err error) {
	db := s.Db.GetGormDb()

	if err := db.Model(&Sim{}).Count(&simsCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	if err := db.Model(&Sim{}).Where("status = ?", SimStatusActive).Count(&activeCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	if err := db.Model(&Sim{}).Where("status = ?", SimStatusInactive).Count(&deactiveCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	if err := db.Model(&Sim{}).Where("status = ?", SimStatusTerminated).Count(&terminatedCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	return simsCount, activeCount, deactiveCount, terminatedCount, nil
}
