package db

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Org struct {
	gorm.Model
	Name  string `gorm:"not null;type:string;uniqueIndex:orgname_idx_case_insensetive,expression:lower(name),where:deleted_at is null"`
	Imsis []Imsi
}

// Represents record in HSS db
type Imsi struct {
	gorm.Model
	OrgID          uint `gorm:"not null"`
	Org            *Org
	Imsi           string `gorm:"index:imsi_unique_idx,unique,where:deleted_at is null;not null;size:15;check:imsi_checker,imsi ~ $$^\\d+$$"` // IMSI Sim ID  (International mobile subscriber identity) https://www.netmanias.com/en/post/blog/5929/lte/lte-user-identifiers-imsi-and-guti
	Op             []byte `gorm:"size:16;"`                                                                                                   // Pre Shared Key. This is optional and configured in operator’s DB in Authentication center and USIM. https://www.3glteinfo.com/lte-security-architecture/
	Amf            []byte `gorm:"size:2;"`                                                                                                    // Pre Shared Key. Configured in operator’s DB in Authentication center and USIM
	Key            []byte `gorm:"size:16;"`                                                                                                   // Key from the SIM
	AlgoType       uint32
	UeDlAmbrBps    uint32
	UeUlAmbrBps    uint32
	Sqn            uint32
	CsgIdPrsent    bool
	CsgId          uint32
	DefaultApnName string
	UserUuid       uuid.UUID `gorm:"not null;type:uuid"`
	Tai            Tai
}

// Tracking Area Identity (TAI)
// Assumption: one IMIS can have only one tracking area
type Tai struct {
	gorm.Model
	ImsiID          uint      `gorm:"uniqueIndex:tai_imsi_unique_idx;not null"`
	PlmId           string    `gorm:"size:6;uniqueIndex:tai_imsi_unique_idx;not null"` // Public Land Mobile Network Identity (MCC+MNC)
	Tac             uint32    `gorm:"uniqueIndex:tai_imsi_unique_idx,where:deleted_at is null;not null"`
	DeviceUpdatedAt time.Time // time when it was updated on the device
}

type User struct {
	gorm.Model
	ID        uint      `gorm:"primarykey"`
	Uuid      uuid.UUID `gorm:"uniqueIndex:uuid_unique,where:deleted_at is null;not null;type:uuid"`
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type Guti struct {
	CreatedAt       time.Time // do not set it directly, it will be overridden
	DeviceUpdatedAt time.Time // time when it was updated on the device
	Imsi            string    `gorm:"uniqueIndex;not null;size:15;check:imsi_checker,imsi ~ $$^\\d+$$"`
	Plmn_id         string    `gorm:"uniqueIndex:idx_guti;not null;size:6"`
	Mmegi           uint32    `gorm:"uniqueIndex:idx_guti;not null"`
	Mmec            uint32    `gorm:"uniqueIndex:idx_guti;not null"`
	MTmsi           uint32    `gorm:"uniqueIndex:idx_guti;not null"`
}
