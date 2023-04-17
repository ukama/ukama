package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Uuid         uuid.UUID `gorm:"primaryKey;type:uuid"`
	OwnerId      uuid.UUID
	Name         string
	SimType      SimType
	OrgId        uuid.UUID       `gorm:"not null;type:uuid;index"`
	Active       bool            `gorm:"not null; default:false"`
	Duration     uint64          `gorm:"not null; default:0"`
	SmsVolume    uint64          `gorm:"not null; default:0"`
	DataVolume   uint64          `gorm:"not null; default:0"`
	VoiceVolume  uint64          `gorm:"not null; default:0"`
	Rate         *PackageRate    `gorm:"not null"`
	Markup       *PackageMarkup  `gorm:"not null"`
	Details      *PackageDetails `gorm:"not null"`
	Type         PackageType     `gorm:"not null; type: uint; default:2"`
	DataUnits    DataUnitType    `gorm:"not null; type: uint; default:3"`
	CallUnits    CallUnitType    `gorm:"not null; type: uint; default:2"`
	MessageUnits MessageUnitType `gorm:"not null; type: uint; default:1"`
	Flatrate     bool            `gorm:"not null; default:false"`
	Currency     string          `gorm:"not null; default:Dollar"`
	EffectiveAt  time.Time       `gorm:"not null"`
	EndAt        time.Time
	Country      string `gorm:"not null;type:string"`
	Provider     string `gorm:"not null;type:string"`
}

type PackageDetails struct {
	gorm.Model
	Dlbr uint64 `gorm:"not null; default:10240000"`
	Ulbr uint64 `gorm:"not null; default:10240000"`
	Apn  string
}

type PackageRate struct {
	gorm.Model
	Amount float64 `gorm:"type:float"`
	SmsMo  float64 `gorm:"type:float"`
	SmsMt  float64 `gorm:"type:float"`
	Data   float64 `gorm:"type:float"`
}

/* View only for owners */
type PackageMarkup struct {
	gorm.Model
	BaseRateId uuid.UUID `gorm:"not null;type:uuid"`
	Markup     float64   `gorm:"type:float"`
}
