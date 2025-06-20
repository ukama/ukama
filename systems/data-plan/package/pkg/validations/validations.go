/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validations

import (
	"errors"
	"time"
)

func IsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func IsReqEmpty(id uint64) bool {
	return id == 0
}

func ValidateDOB(dob string) (time.Time, error) {
	t, err := time.Parse("02-01-2006", dob)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, must be dd-mm-yyyy")
	}
	return t, nil
}
