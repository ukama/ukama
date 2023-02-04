package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type Subscriber struct {
	SubscriberID          uuid.UUID `gorm:"primaryKey;type:uuid;unique"`
	FirstName             string    `gorm:"size:255"`
	LastName              string    `gorm:"size:255"`
	NetworkID             uuid.UUID `gorm:"type:uuid;index"`
<<<<<<< HEAD
	OrgID          uuid.UUID `gorm:"type:uuid"`
=======
	OrgID                 uuid.UUID `gorm:"type:uuid"`
>>>>>>> subscriber-sys_sim-manager
	Email                 string    `gorm:"size:255"`
	PhoneNumber           string    `gorm:"size:15"`
	Gender                string    `gorm:"size:255"`
	DOB                   time.Time
	ProofOfIdentification string `gorm:"size:255"`
	IdSerial              string `gorm:"size:255"`
	Address               string `gorm:"size:255"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time `sql:"index"`
}
