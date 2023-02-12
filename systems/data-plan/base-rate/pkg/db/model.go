package db

import (
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type Rate struct {
	gorm.Model
	RateID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Country      string
	Network      string
	Vpmn         string
	Imsi         string
	SmsMo       string
	SmsMt       string
	Data         string
	X2g          string
	X3g          string
	X5g          string
	Lte          string
	LteM        string
	Apn          string
	EffectiveAt string
	EndAt       string
	SimType     string
}
