package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type MailerRepo interface {
	SendEmail(mail *Mailing) error
	GetEmailById(mailerId uuid.UUID) (*Mailing, error)
}

type mailerRepo struct {
	Db sql.Db
}

func NewMailerRepo(db sql.Db) MailerRepo {
	return &mailerRepo{
		Db: db,
	}
}

func (s *mailerRepo) SendEmail(mail *Mailing) error {
	db := s.Db.GetGormDb()
	return db.Create(mail).Error
}

func (s *mailerRepo) GetEmailById(mailerId uuid.UUID) (*Mailing, error) {
	db := s.Db.GetGormDb()
	mail := &Mailing{}
	err := db.Where("mail_id = ?", mailerId).First(mail).Error
	return mail, err
}
