/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validation

import (
	"fmt"
	"time"
)

func IsFutureDate(date string) error {
	t, err := FromString(date)
	if err != nil {
		return fmt.Errorf("invalid date format, must be RFC3339 standard: %w", err)
	}

	if t.After(time.Now()) {
		return nil
	}

	return fmt.Errorf("date %s is not in the future", t.Format(time.RFC3339))
}

func IsAfterDate(date string, after string) error {
	t, err := FromString(date)
	if err != nil {
		return fmt.Errorf("invalid date format, must be RFC3339 standard: %w", err)
	}

	a, err := FromString(after)
	if err != nil {
		return fmt.Errorf("invalid date format, must be RFC3339 standard: %w", err)
	}

	if t.After(a) {
		return nil
	}

	return fmt.Errorf("date is not after %s", a.Format(time.RFC3339))
}

func ValidateDate(date string) (string, error) {
	t, err := FromString(date)
	if err != nil {
		return "", fmt.Errorf("invalid date format, must be RFC3339 standard: %w", err)
	}

	return t.Format(time.RFC3339), nil
}

func FromString(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
