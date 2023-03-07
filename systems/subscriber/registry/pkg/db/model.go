package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Subscriber struct {
	SubscriberID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	FirstName             string    `gorm:"size:255"`
	LastName              string    `gorm:"size:255"`
	NetworkID             uuid.UUID `gorm:"type:uuid;index"`
	OrgID                 uuid.UUID `gorm:"type:uuid"`
	Email                 string    `gorm:"size:255"`
	PhoneNumber           string    `gorm:"size:15"`
	Gender                string    `gorm:"size:255"`
	DOB                   string
	ProofOfIdentification string `gorm:"size:255"`
	IdSerial              string `gorm:"size:255"`
	Address               string `gorm:"size:255"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time `sql:"index"`
}
