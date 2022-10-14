package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type simcardRepo struct {
	db sql.Db
}

type SimcardRepo interface {
	Add(simcard *Simcard) error
	Delete(iccid string) error
	DeleteByUser(userUuid uuid.UUID) error
	Get(iccid string) (*Simcard, error)
	UpdateServices(ukama *Service, carrier *Service, nested ...func() error) error
}

func NewSimcardRepo(db sql.Db) *simcardRepo {
	return &simcardRepo{
		db: db,
	}
}

func (s *simcardRepo) Add(sim *Simcard) error {
	d := s.db.GetGormDb().Create(sim)

	return d.Error
}

func (s *simcardRepo) Delete(iccid string) error {
	d := s.db.GetGormDb().Delete(&Simcard{Iccid: iccid})

	return d.Error
}

func (s *simcardRepo) DeleteByUser(userUuid uuid.UUID) error {
	d := s.db.GetGormDb().Exec("delete from simcards where user_id in ( select id from users where uuid  = ?)", userUuid)

	return d.Error
}

func (s *simcardRepo) Get(iccid string) (*Simcard, error) {
	var sim Simcard
	d := s.db.GetGormDb().Preload(clause.Associations).Where("iccid = ?", iccid).First(&sim)

	return &sim, d.Error
}

// UpdateServices updates services for simcard
// nested func used to execute code in scope of transaction
// in nested returns error then transaction is rolled back
func (s *simcardRepo) UpdateServices(ukama *Service, carrier *Service, nested ...func() error) error {
	err := s.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		d := tx.Save(ukama)
		if d.Error != nil {
			return d.Error
		}
		d = tx.Save(carrier)
		if d.Error != nil {
			return d.Error
		}

		for _, f := range nested {
			err := f()
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
