package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Uuid           uuid.UUID `gorm:"unique;type:uuid;index"`
	OwnerId        uuid.UUID
	Name           string
	SimType        SimType
	OrgId          uuid.UUID      `gorm:"not null;type:uuid;index"`
	Active         bool           `gorm:"not null; default:false"`
	Duration       uint64         `gorm:"not null; default:0"`
	SmsVolume      uint64         `gorm:"not null; default:0"`
	DataVolume     uint64         `gorm:"not null; default:0"`
	VoiceVolume    uint64         `gorm:"not null; default:0"`
	PackageRate    PackageRate    `gorm:"foreignKey:PackageID;references:Uuid"`
	PackageMarkup  PackageMarkup  `gorm:"foreignKey:PackageID;references:Uuid"`
	PackageDetails PackageDetails `gorm:"foreignKey:PackageID;references:Uuid"`
	Type           ukama.PackageType
	DataUnits      ukama.DataUnitType
	VoiceUnits     ukama.CallUnitType
	MessageUnits   ukama.MessageUnitType
	Flatrate       bool      `gorm:"not null; default:false"`
	Currency       string    `gorm:"not null; default:Dollar"`
	From           time.Time `gorm:"not null"`
	To             time.Time `gorm:"not null"`
	Country        string    `gorm:"not null;type:string"`
	Provider       string    `gorm:"not null;type:string"`
}

type PackageDetails struct {
	gorm.Model
	PackageID uuid.UUID
	Dlbr      uint64 `gorm:"not null; default:10240000"`
	Ulbr      uint64 `gorm:"not null; default:10240000"`
	Apn       string
}

type PackageRate struct {
	gorm.Model
	PackageID uuid.UUID
	Amount    float64 `gorm:"type:float"`
	SmsMo     float64 `gorm:"type:float"`
	SmsMt     float64 `gorm:"type:float"`
	Data      float64 `gorm:"type:float"`
}

/* View only for owners */
type PackageMarkup struct {
	gorm.Model
	PackageID  uuid.UUID
	BaseRateId uuid.UUID `gorm:"not null;type:uuid"`
	Markup     float64   `gorm:"type:float"`
}
