package db

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	SubscriberID          uuid.UUID `gorm:"type:uuid"`
	FirstName             string
	LastName              string
	NetworkID          uuid.UUID `gorm:"type:uuid"`
	Email                 string
	PhoneNumber           string
	Gender                string
	DOB       time.Time `gorm:"column:dob"`
	ProofOfIdentification string
	IdSerial              string
	Address               string
}
