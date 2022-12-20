package db

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Subscriber struct {
	gorm.Model
	SubscriberID          uuid.UUID `gorm:"type:uuid"`
	FirstName 				  string
	LastName                   string
	SimID				  string
	Email                 string
	PhoneNumber           string
	Gender				string
	DOB                   string
	ProofOfIdentification string
	IdSerial              string
	Address               string
}



