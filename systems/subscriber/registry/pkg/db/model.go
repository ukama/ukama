/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)
 
 type Subscriber struct {
	 SubscriberId          uuid.UUID `gorm:"primaryKey;type:uuid"`
	 Name             string    `gorm:"size:255"`
	 NetworkId             uuid.UUID `gorm:"type:uuid;index"`
	 Email                 string    `gorm:"size:255;not null;index"` 
	 PhoneNumber           string    `gorm:"size:15"`
	 Gender                string    `gorm:"size:255"`
	 DOB                   string
	 ProofOfIdentification string `gorm:"size:255"`
	 IdSerial              string `gorm:"size:255"`
	 Address               string `gorm:"size:255"`
	 CreatedAt             time.Time
	 UpdatedAt             time.Time
	 DeletedAt             *time.Time `sql:"index"`
 }