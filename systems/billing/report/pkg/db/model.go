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
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
)

type Report struct {
	Id        uuid.UUID      `gorm:"primaryKey;type:uuid"`
	OwnerId   uuid.UUID      `gorm:"uniqueIndex:report_id_period,where:deleted_at is null;not null;type:uuid"`
	OwnerType OwnerType      `gorm:"not null"`
	NetworkId uuid.UUID      `gorm:"type:uuid"`
	Type      ReportType     `gorm:"not null"`
	Period    time.Time      `gorm:"uniqueIndex:report_id_period,where:deleted_at is null;not null"`
	RawReport datatypes.JSON `gorm:"not null"`
	IsPaid    bool           `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type OwnerType uint8

const (
	OwnerTypeUnknown OwnerType = iota
	OwnerTypeOrg
	OwnerTypeSubscriber
)

func (s *OwnerType) Scan(value interface{}) error {
	*s = OwnerType(uint8(value.(int64)))

	return nil
}

func (o OwnerType) Value() (driver.Value, error) {
	return int64(o), nil
}

func (o OwnerType) String() string {
	t := map[OwnerType]string{0: "unknown", 1: "org", 2: "subscriber"}

	v, ok := t[o]
	if !ok {
		return t[0]
	}

	return v
}

func ParseOwnerType(value string) OwnerType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return OwnerType(i)
	}

	t := map[string]OwnerType{"unknown": 0, "org": 1, "subscriber": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return OwnerType(0)
	}

	return OwnerType(v)
}

type ReportType uint8

const (
	ReportTypeUnknown ReportType = iota
	ReportTypeInvoice
	ReportTypeConsumption
)

func (r *ReportType) Scan(value interface{}) error {
	*r = ReportType(uint8(value.(int64)))

	return nil
}

func (r ReportType) Value() (driver.Value, error) {
	return int64(r), nil
}

func (r ReportType) String() string {
	t := map[ReportType]string{0: "unknown", 1: "invoice", 2: "consumption"}

	v, ok := t[r]
	if !ok {
		return t[0]
	}

	return v
}

func ParseReportType(value string) ReportType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return ReportType(i)
	}

	t := map[string]ReportType{"unknown": 0, "invoice": 1, "consumption": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return ReportType(0)
	}

	return ReportType(v)
}
