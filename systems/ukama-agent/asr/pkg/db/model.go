/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type StatusReason int64

const (
	UNKNOWN StatusReason = iota
	ACTIVATION
	PACKAGE_UPDATE
	DEACTIVATION
	NO_DATA_AVAILABLE
	POLICY_FAILURE
)

// Represents record in HSS db
type Asr struct {
	gorm.Model

	//(TODO: will fix these check constaints later) Iccid string `gorm:"index:asr_iccid_idx,unique,where:deleted_at is null;not null;size:22;check:iccid_checker,iccid ~ $$^\\d+$$"`
	Iccid string `gorm:"index:asr_iccid_idx,unique,where:deleted_at is null;not null;size:22"`
	//IMSI might not be unique as same IMSI might be authorized to use multiple network of Org which means multiple enetry for the IMSI in HLR or may be use many to many relattion here.
	// For Now we are considering that use case where each imsi could only belong to one network.
	// IMSI Sim ID  (International mobile subscriber identity) https://www.netmanias.com/en/post/blog/5929/lte/lte-user-identifiers-imsi-and-guti
	//(TODO: will fix these check constaints later) Imsi string `gorm:"index:asr_imsi_idx,unique,where:deleted_at is null;not null;size:15;check:asr_checker,imsi ~ $$^\\d+$$"`
	Imsi string `gorm:"index:asr_imsi_idx,unique,where:deleted_at is null;not null;size:15"`
	// Pre Shared Key. This is optional and configured in operator’s DB in Authentication center and USIM. https://www.3glteinfo.com/lte-security-architecture/
	Op []byte `gorm:"size:16;"`
	// Pre Shared Key. Configured in operator’s DB in Authentication center and USIM
	Amf []byte `gorm:"size:2;"`
	// Key from the SIM
	Key                     []byte `gorm:"size:16;"`
	AlgoType                uint32
	UeDlAmbrBps             uint32
	UeUlAmbrBps             uint32
	Sqn                     uint64
	CsgIdPrsent             bool
	CsgId                   uint32
	DefaultApnName          string
	NetworkId               uuid.UUID `gorm:"not null;type:uuid"`
	Tai                     Tai
	PackageId               uuid.UUID `gorm:"not null;type uuid"`
	Policy                  Policy
	LastStatusChangeAt      time.Time
	AllowedTimeOfService    int64
	LastStatusChangeReasons StatusReason
}

// Tracking Area Identity (TAI)
// Assumption: one IMIS can have only one tracking area
type Tai struct {
	gorm.Model
	AsrID           uint      `gorm:"uniqueIndex:tai_asr_unique_idx;not null"`
	PlmnId          string    `gorm:"size:6;uniqueIndex:tai_asr_unique_idx;not null"` // Public Land Mobile Network Identity (MCC+MNC)
	Tac             uint32    `gorm:"uniqueIndex:tai_asr_unique_idx,where:deleted_at is null;not null"`
	DeviceUpdatedAt time.Time // time when it was updated on the device
}

type Guti struct {
	CreatedAt       time.Time // do not set it directly, it will be overridden
	DeviceUpdatedAt time.Time // time when it was updated on the device
	//(TODO: will fix these check constaints later) Imsi            string    `gorm:"primarykey;uniqueIndex:guti_asr_unique_idx;not null;size:15;check:imsi_checker,imsi ~ $$^\\d+$$"`
	Imsi   string `gorm:"primarykey;uniqueIndex:guti_asr_unique_idx;not null;size:15"`
	PlmnId string `gorm:"uniqueIndex:guti_asr_unique_idx;not null;size:6"`
	Mmegi  uint32 `gorm:"uniqueIndex:guti_asr_unique_idx;not null"`
	Mmec   uint32 `gorm:"uniqueIndex:guti_asr_unique_idx;not null"`
	MTmsi  uint32 `gorm:"uniqueIndex:guti_asr_unique_idx;not null"`
}

type Policy struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Id           uuid.UUID      `gorm:"primarykey;type:uuid"`
	Burst        uint64
	TotalData    uint64
	ConsumedData uint64
	Dlbr         uint64
	Ulbr         uint64
	StartTime    uint64
	EndTime      uint64
	AsrID        uint
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
	case "POLICY_FAILURE", "policy_failure":
		return StatusReason(POLICY_FAILURE)
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
	case POLICY_FAILURE:
		return "POLICY_FAILURE"
	default:
		return "UNKNOWN"
	}
}

func (s StatusReason) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s *StatusReason) Scan(value interface{}) error {
	val, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid status value %v %T", value, value)
	}

	switch StatusReason(val) {
	case ACTIVATION, PACKAGE_UPDATE, DEACTIVATION, NO_DATA_AVAILABLE, POLICY_FAILURE:
		*s = StatusReason(val)
	default:
		*s = UNKNOWN
	}
	return nil
}
