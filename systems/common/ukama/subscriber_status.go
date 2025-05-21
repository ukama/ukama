/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama

import (
	"database/sql/driver"
	"strconv"
	"strings"
)


type SubscriberStatus uint8

const (
	SubscriberStatusUnknown        SubscriberStatus = iota 
	SubscriberStatusActive                               
	SubscriberStatusPendingDeletion                       
)

func (s *SubscriberStatus) Scan(value interface{}) error {
	*s = SubscriberStatus(uint8(value.(int64)))
	return nil
}

func (s SubscriberStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s SubscriberStatus) String() string {
	t := map[SubscriberStatus]string{
		0: "unknown", 
		1: "active", 
		2: "pending_deletion",
	}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseSubscriberStatus(value string) SubscriberStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return SubscriberStatus(i)
	}

	t := map[string]SubscriberStatus{
		"unknown":          0, 
		"active":           1, 
		"pending_deletion": 2,
	}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return SubscriberStatus(0)
	}

	return SubscriberStatus(v)
}