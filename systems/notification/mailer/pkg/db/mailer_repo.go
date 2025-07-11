/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type MailerRepo interface {
	CreateEmail(mail *Mailing) error
	GetEmailById(mailerId uuid.UUID) (*Mailing, error)
	UpdateEmailStatus(mailing *Mailing) error
	GetFailedEmails() ([]*Mailing, error)
}

type mailerRepo struct {
	Db sql.Db
}

func NewMailerRepo(db sql.Db) MailerRepo {
	return &mailerRepo{
		Db: db,
	}
}

func (s *mailerRepo) CreateEmail(mail *Mailing) error {
	db := s.Db.GetGormDb()
	return db.Create(mail).Error
}

func (s *mailerRepo) GetEmailById(mailerId uuid.UUID) (*Mailing, error) {
	db := s.Db.GetGormDb()
	mail := &Mailing{}
	err := db.Where("mail_id = ?", mailerId).First(mail).Error
	return mail, err
}
func (r *mailerRepo) UpdateEmailStatus(mailing *Mailing) error {
	db := r.Db.GetGormDb()

	result := db.Model(&Mailing{}).
		Where("mail_id = ?", mailing.MailId).
		Updates(map[string]interface{}{
			"status":          mailing.Status,
			"retry_count":     mailing.RetryCount,
			"next_retry_time": mailing.NextRetryTime,
			"updated_at":      time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *mailerRepo) GetFailedEmails() ([]*Mailing, error) {
	var mailings []*Mailing
	db := r.Db.GetGormDb()

	result := db.Where("status IN (?, ?) AND retry_count < ? AND (next_retry_time <= ? OR next_retry_time IS NULL)",
		ukama.MailStatusFailed,
		ukama.MailStatusRetry,
		ukama.MaxRetryCount,
		time.Now(),
	).Find(&mailings)

	return mailings, result.Error
}
