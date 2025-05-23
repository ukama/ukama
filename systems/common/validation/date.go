/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validation

import (
	"errors"
	"time"
)

func IsFutureDate(date string) error {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return err
	}
	if t.After(time.Now()) {
		return nil
	}
	return errors.New("date is not in the future")
}

func ValidateDate(date string) (string, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", errors.New("invalid date format, must be RFC3339 standard")
	}
	return t.Format(time.RFC3339), nil
}

func IsAfterDate(date string, after string) error {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return errors.New("invalid date format, must be RFC3339 standard")
	}

	a, err := time.Parse(time.RFC3339, after)
	if err != nil {
		return errors.New("invalid date format, must be RFC3339 standard")
	}
	if t.After(a) {
		return nil
	}
	return errors.New("date is not after" + a.Format(time.RFC3339))
}

func FromString(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
