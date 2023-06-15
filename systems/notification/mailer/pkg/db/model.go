package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Mailing struct {
	MailId    uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Email     string     `gorm:"not null"`
	Subject   string     `gorm:"not null"`
	Body      string     `gorm:"not null"`
	SentAt    *time.Time `gorm:"index"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
}
