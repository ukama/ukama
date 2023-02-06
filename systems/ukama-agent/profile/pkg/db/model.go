package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type PackageDetails struct {
	PackageId            uuid.UUID
	AllowedTimeOfService time.Duration
	TotalDataBytes       uint64
	ConsumedDataBytes    uint64
}

// Represents record in HSS db
type Profile struct {
	gorm.Model

	Iccid string `gorm:"index:asr_iccid_idx,unique,where:deleted_at is null;not null;size:22;check:iccid_checker,iccid ~ $$^\\d+$$"`
	//IMSI might not be unique as same IMSI might be authorized to use multiple network of Org which means multiple enetry for the IMSI in HLR or may be use many to many relattion here.
	// For Now we are considering that use case where each imsi could only belong to one network.
	// IMSI Sim ID  (International mobile subscriber identity) https://www.netmanias.com/en/post/blog/5929/lte/lte-user-identifiers-imsi-and-guti
	Imsi string `gorm:"index:asr_imsi_idx,unique,where:deleted_at is null;not null;size:15;check:asr_checker,imsi ~ $$^\\d+$$"`
	// Pre Shared Key. This is optional and configured in operator’s DB in Authentication center and USIM. https://www.3glteinfo.com/lte-security-architecture/

	UeDlBps              uint64        `gorm:"default:100000"` //TODO: Add it to Package DB in data-plan
	UeUlBps              uint64        `gorm:"default:10000"`  //TODO: Add it to Package DB in data-plan
	ApnName              string        `gorm:"default:ukama"`  //TODO: Add it to Package DB in data-plan
	NetworkID            uuid.UUID     `gorm:"not null;type:uuid"`
	PackageId            uuid.UUID     `gorm:"not null;type uuid"`
	AllowedTimeOfService time.Duration `gorm:"default:43200s"` //TODO: Add it to Package DB in data-plan (30*24*60=43200)
	TotalDataBytes       uint64
	ConsumedDataBytes    uint64
}

/* May be ignore metrics here and use CDR fir that */
type Metrics struct {
	gorm.Model
	ProfileID             uint `gorm:"uniqueIndex:metrics_unique_idx;not null"`
	LastSessionAt         time.Time
	LastSessionTimePeriod time.Duration
	LastUsedBytes         uint64
}
