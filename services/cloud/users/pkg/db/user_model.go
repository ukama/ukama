package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Org struct {
	gorm.Model
	Name      string `gorm:"not null;type:string;uniqueIndex:orgname_idx_case_insensetive,expression:lower(name),where:deleted_at is null"`
	UserLimit uint
	Users     []User
}

type User struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	Name        string    `gorm:"not null;default:'unknown'"`
	Email       string
	Phone       string
	Simcard     Simcard
	OrgID       uint `gorm:"not null;default:1"`
	Org         *Org
	Deactivated bool
}

// storage for service statuses
type Service struct {
	gorm.Model
	Iccid string `gorm:"uniqueIndex:iccid_type_unique,where:deleted_at is null"`
	Voice bool
	Sms   bool
	Data  bool
	Type  uint8 `gorm:"uniqueIndex:iccid_type_unique,where:deleted_at is null"` // 0 - ukama, 1 - carrier
}

const (
	ServiceTypeUkama   = 0
	ServiceTypeCarrier = 1
)

type Simcard struct {
	UserID     uint       `gorm:"not null;uniqueIndex;default:0"`
	IsPhysical bool       `gorm:"not null;default:false"`
	Iccid      string     `gorm:"primarykey"`
	Services   []*Service `gorm:"foreignKey:Iccid;references:Iccid"`
	Source     string
}

func (s Simcard) GetServices(srvType uint8) *Service {
	for _, sr := range s.Services {
		if sr.Type == srvType {
			return sr
		}
	}

	return nil
}
