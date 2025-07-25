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
	"time"

	"github.com/ukama/ukama/systems/common/ukama"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type JSONMap map[string]interface{}

type Mailing struct {
	MailId        uuid.UUID        `gorm:"primaryKey;type:uuid"`
	Email         string           `gorm:"size:255"`
	TemplateName  string           `gorm:"size:255"`
	SentAt        *time.Time       `gorm:"index"`
	Status        ukama.MailStatus `gorm:"type:uint;not null"`
	RetryCount    int              `gorm:"default:0"`
	NextRetryTime *time.Time       `gorm:"index"`
	Values        JSONMap          `gorm:"type:jsonb" json:"values"`
	CreatedAt     time.Time        `gorm:"not null"`
	UpdatedAt     time.Time        `gorm:"not null"`
	DeletedAt     *time.Time       `gorm:"index"`
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
