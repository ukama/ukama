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

type MailStatus uint8

// TODO: we need unknown as a sentinel value for all non valid values of the universe. We cannot
// rely on pending (which is a valid status)
const (
	MailStatusPending MailStatus = iota
	MailStatusSuccess
	MailStatusFailed
	MailStatusRetry
	MailStatusProcess

	MaxRetryCount = 3
)

func (s *MailStatus) Scan(value interface{}) error {
	*s = MailStatus(uint8(value.(int64)))

	return nil
}

func (s MailStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s MailStatus) String() string {
	t := map[MailStatus]string{0: "pending", 1: "success", 2: "failed", 3: "retry", 4: "process"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParseMailStatus(value string) MailStatus {
	i, err := strconv.Atoi(value)
	if err == nil {
		return MailStatus(i)
	}

	t := map[string]MailStatus{"pending": 0, "success": 1, "failed": 2, "retry": 3, "process": 4}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return MailStatus(0)
	}

	return MailStatus(v)
}
