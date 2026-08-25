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

type PolicyViolationReason uint8

const (
	PolicyViolationReasonUnknown PolicyViolationReason = iota
	PolicyViolationReasonDataCapExceeded
	PolicyViolationReasonPackageExpired
)

func (s *PolicyViolationReason) Scan(value interface{}) error {
	*s = PolicyViolationReason(uint8(value.(int64)))

	return nil
}

func (s PolicyViolationReason) Value() (driver.Value, error) {
	return int64(s), nil
}

func (s PolicyViolationReason) String() string {
	t := map[PolicyViolationReason]string{0: "unknown", 1: "data_cap_exceeded", 2: "package_expired"}

	v, ok := t[s]
	if !ok {
		return t[0]
	}

	return v
}

func ParsePolicyViolationReason(value string) PolicyViolationReason {
	i, err := strconv.Atoi(value)
	if err == nil {
		return PolicyViolationReason(i)
	}

	t := map[string]PolicyViolationReason{"unknown": 0, "data_cap_exceeded": 1, "package_expired": 2}

	v, ok := t[strings.ToLower(value)]
	if !ok {
		return PolicyViolationReason(0)
	}

	return v
}
