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
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Org struct {
	Id          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"uniqueIndex"`
	Deactivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Network struct {
	Id               uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name             string    `gorm:"uniqueIndex:network_name_org_idx"`
	OrgId            uuid.UUID `gorm:"uniqueIndex:network_name_org_idx;type:uuid"`
	Org              *Org
	Deactivated      bool
	AllowedCountries pq.StringArray `gorm:"type:varchar(64)[]" json:"allowed_countries"`
	AllowedNetworks  pq.StringArray `gorm:"type:varchar(64)[]" json:"allowed_networks"`
	Budget           float64
	Overdraft        float64
	TrafficPolicy    uint32
	PaymentLinks     bool
	Country          string `json:"country"`
	Language         LanguageType
	Currency         string `json:"currency"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	SyncStatus       ukama.StatusType
}

type LanguageType uint8

const (
	UnknownLanguage LanguageType = iota
	EN                           = 1
	FR                           = 2
)

func (l *LanguageType) Scan(value interface{}) error {
	*l = LanguageType(uint8(value.(int64)))
	return nil
}

func (l LanguageType) Value() (driver.Value, error) {
	return int64(l), nil
}

func (l LanguageType) String() string {
	t := map[LanguageType]string{0: "unknown", 1: "fr", 2: "en"}

	v, ok := t[l]
	if !ok {
		return t[0]
	}

	return v
}

func ParseType(value string) LanguageType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return LanguageType(i)
	}

	t := map[string]LanguageType{"unknown": 0, "fr": 1, "en": 2}

	v, ok := t[value]
	if !ok {
		return LanguageType(0)
	}

	return LanguageType(v)
}
