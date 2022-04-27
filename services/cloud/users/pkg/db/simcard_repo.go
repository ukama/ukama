package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/services/common/sql"
)

type simcardRepo struct {
	db sql.Db
}

type SimcardRepo interface {
	Add(simcard *Simcard) error
	Delete(iccid string) error
	DeleteByUser(userUuid uuid.UUID) error
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
