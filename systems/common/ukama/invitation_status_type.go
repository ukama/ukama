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

type InvitationStatusType uint8

const (
	Pending  InvitationStatusType = 0
	Accepted InvitationStatusType = 1
	Declined InvitationStatusType = 2
)

func (s *InvitationStatusType) Scan(value interface{}) error {
	*s = InvitationStatusType(uint8(value.(int64)))

	return nil
}

func (s InvitationStatusType) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s InvitationStatusType) String() string {
	t := map[InvitationStatusType]string{0: "pending", 1: "accepted", 2: "declined"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseInvitationStatusType(value string) InvitationStatusType {
	i, err := strconv.Atoi(value)
	if err == nil {
		return InvitationStatusType(i)
	}

	t := map[string]InvitationStatusType{"pending": 0, "accepted": 1, "declined": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return InvitationStatusType(0)
	}

	return InvitationStatusType(v)
}
