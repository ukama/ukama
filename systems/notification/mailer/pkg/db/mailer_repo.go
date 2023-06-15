package db

import (
	"github.com/ukama/ukama/systems/common/sql"
)

type MaillingRepo interface {
	SendEmail(mail *Mailing) error
}

type maillingRepo struct {
	Db sql.Db
}

func NewMaillingRepo(db sql.Db) MaillingRepo {
	return &maillingRepo{
		Db: db,
	}
}

func (s *maillingRepo) SendEmail(mail *Mailing) error {
	db := s.Db.GetGormDb()
	return db.Create(mail).Error
}