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

// TODO: should be renamed as MailStatus instead of just status + update the new name
// on all other dependent services that are using the current name: status to MailStatus
type Status uint8

// TODO: we need unknown as a sentinel value for all non valid values of the universe. We cannot
// rely on pending (which is a valid status)
const (
	Pending Status = iota
	Success
	Failed
	Retry
	Process

	MaxRetryCount = 3
)

func (s *Status) Scan(value interface{}) error {
	*s = Status(uint8(value.(int64)))

	return nil
}

func (s Status) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s Status) String() string {
	t := map[Status]string{0: "pending", 1: "success", 2: "failed", 3: "retry", 4: "process"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

// TODO: should be renamed as ParseMailStatus + make all necessary updates on other dependent services
// that are using this.
func ParseStatus(value string) Status {
	i, err := strconv.Atoi(value)
	if err == nil {
		return Status(i)
	}

	t := map[string]Status{"pending": 0, "success": 1, "failed": 2, "retry": 3, "process": 4}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return Status(0)
	}

	return Status(v)
}
