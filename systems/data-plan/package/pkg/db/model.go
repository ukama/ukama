package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	Uuid         uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name         string
	SimType      SimType
	OrgId        uuid.UUID       `gorm:"not null;type:uuid;index"`
	Active       bool            `gorm:"not null; default:false"`
	Duration     uint            `gorm:"not null; default:0"`
	SmsVolume    uint            `gorm:"not null; default:0"`
	DataVolume   uint            `gorm:"not null; default:0"`
	VoiceVolume  uint            `gorm:"not null; default:0"`
	Dlbr         uint64          `gorm:"not null; default:0"`
	Ulbr         uint64          `gorm:"not null; default:0"`
	Rate         *Rate           `gorm:"not null"`
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
	Apn          string
	X2g          bool `gorm:"type:bool; default:false"`
	X3g          bool `gorm:"type:bool; default:false"`
	X5g          bool `gorm:"type:bool; default:false"`
	Lte          bool `gorm:"type:bool; default:false"`
	LteM         bool `gorm:"type:bool; default:false"`
}

type Rate struct {
	gorm.Model
	BaseRateId uuid.UUID `gorm:"not null;type:uuid"`
	Markup     float64   `gorm:"type:float"`
	Flatrate   float64   `gorm:"type:float"`
	Vpmn       string
	Imsi       int64
	SmsMo      float64 `gorm:"type:float"`
	SmsMt      float64 `gorm:"type:float"`
	Data       float64 `gorm:"type:float"`
}
