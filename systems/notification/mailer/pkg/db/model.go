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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type JSONMap map[string]interface{}
type Status uint8

const (
	Pending Status = iota
	Success
	Failed
	Retry
	Process
	MaxRetryCount = 3
)

type Mailing struct {
	MailId        uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Email         string     `gorm:"size:255"`
	TemplateName  string     `gorm:"size:255"`
	SentAt        *time.Time `gorm:"index"`
	Status        Status     `gorm:"type:uint;not null"`
	RetryCount    int        `gorm:"default:0"`
	NextRetryTime *time.Time `gorm:"index"`
	Values        JSONMap    `gorm:"type:jsonb" json:"values"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     time.Time  `gorm:"not null"`
	DeletedAt     *time.Time `gorm:"index"`
}

func (m JSONMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(JSONMap)
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONMap: %v", value)
	}
	return json.Unmarshal(data, m)
}

func (s *Status) Scan(value interface{}) error {
	*s = Status(uint8(value.(int64)))

	return nil
}

func (s Status) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s Status) String() string {
	t := map[Status]string{0: "pending", 1: "success", 2: "Failed", 3: "retry", 4: "process"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseStatus(value string) Status {
	i, err := strconv.Atoi(value)
	if err == nil {
		return Status(i)
	}

	t := map[string]Status{"unknown": 0, "onboarded": 1, "configured": 2, "operational": 3, "faulty": 4}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return Status(0)
	}

	return Status(v)
}
