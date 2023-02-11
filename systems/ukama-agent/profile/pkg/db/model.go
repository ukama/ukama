package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type StatusReason int

const (
	UNKNOWN StatusReason = iota
	ACTIVATION
	PACKAGE_UPDATE
	DEACTIVATION
	NO_DATA_AVAILABLE
)

type PackageDetails struct {
	PackageId            uuid.UUID
	AllowedTimeOfService time.Duration
	TotalDataBytes       uint64
	ConsumedDataBytes    uint64
	UeDlBps              uint64
	UeUlBps              uint64
	ApnName              string
	LastStatusChangeAt   time.Time
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

	UeDlBps                 uint64    `gorm:"default:100000"` //TODO: Add it to Package DB in data-plan
	UeUlBps                 uint64    `gorm:"default:10000"`  //TODO: Add it to Package DB in data-plan
	ApnName                 string    `gorm:"default:ukama"`  //TODO: Add it to Package DB in data-plan
	NetworkId               uuid.UUID `gorm:"not null;type:uuid"`
	PackageId               uuid.UUID `gorm:"not null;type uuid"`
	AllowedTimeOfService    int64     `gorm:"default:43200"` //TODO: Add it to Package DB in data-plan (30*24*60=43200)
	TotalDataBytes          uint64
	ConsumedDataBytes       uint64
	LastStatusChangeAt      time.Time
	LastStatusChangeReasons StatusReason `gorm:"type:int"` // Hold the reason for last status change which is for activation or deactivation
}

/* May be ignore metrics here and use CDR fir that */
type Metrics struct {
	gorm.Model
	ProfileID             uint `gorm:"uniqueIndex:metrics_unique_idx;not null"`
	LastSessionAt         time.Time
	LastSessionTimePeriod time.Duration
	LastUsedBytes         uint64
}

func StatusReasonFromString(s string) StatusReason {
	switch s {
	case "ACTIVATION", "activation":
		return StatusReason(ACTIVATION)
	case "PACKAGE_UPDATE", "package_update":
		return StatusReason(PACKAGE_UPDATE)
	case "DEACTIVATION", "deactivation":
		return StatusReason(DEACTIVATION)
	case "NO_DATA_AVAILABLE", "no_data_available":
		return StatusReason(NO_DATA_AVAILABLE)
	default:
		return StatusReason(UNKNOWN)
	}
}

func (s StatusReason) String() string {
	switch s {
	case ACTIVATION:
		return "ACTIVATION"
	case PACKAGE_UPDATE:
		return "PACKAGE_UPDATE"
	case DEACTIVATION:
		return "DEACTIVATION"
	case NO_DATA_AVAILABLE:
		return "NO_DATA_AVAILABLE"
	default:
		return "UNKNOWN"
	}
}

func (s StatusReason) Value() (driver.Value, error) {
	val := int(s)

	return val, nil
}

func (s *StatusReason) Scan(value interface{}) error {
	val, ok := value.(int)
	if !ok {
		return fmt.Errorf("invalid status value %v", value)
	}
	switch StatusReason(val) {
	case ACTIVATION, PACKAGE_UPDATE, DEACTIVATION, NO_DATA_AVAILABLE:
		*s = StatusReason(val)
	default:
		*s = UNKNOWN
	}

	return fmt.Errorf("invalid status value %v", value)
}
