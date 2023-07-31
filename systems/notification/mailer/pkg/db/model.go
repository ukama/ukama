package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Mailing struct {
	MailId       uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Email        string         `gorm:"size:255"`
	TemplateName string         `gorm:"size:255"`
	SentAt       *time.Time `gorm:"index"`
	Status       string     `gorm:"not null"`
	CreatedAt    time.Time  `gorm:"not null"`
	UpdatedAt    time.Time  `gorm:"not null"`
	DeletedAt    *time.Time `gorm:"index"`
}

