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
	"fmt"
	"time"
)

const (
	MinutesInDay   = 24 * 60
	MinutesInYear  = 365 * MinutesInDay
	MinutesInMonth = 30 * MinutesInDay

	// Safety limits
	MaxPackageYears   = 1000
	MaxAllowedMinutes = uint64(MaxPackageYears * MinutesInYear)
)

var (
	ErrDurationTooLong = errors.New("package duration exceeds the maximum allowed limit of 1000 years")
	ErrDurationZero    = errors.New("package duration must be greater than zero minutes")
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

// ValidatePackageDuration ensures the duration fits safely limits
func ValidatePackageDuration(durationMinutes uint64) error {
	if durationMinutes == 0 {
		return ErrDurationZero
	}

	if durationMinutes > MaxAllowedMinutes {
		return fmt.Errorf("%w: received %d minutes", ErrDurationTooLong, durationMinutes)
	}

	return nil
}

// CalculateEndDate ensures the duration safely converts in a future date
func CalculateEndDate(startDate time.Time, durationMinutes uint64) time.Time {
	years := durationMinutes / MinutesInYear
	remMinutes := durationMinutes % MinutesInYear

	months := remMinutes / MinutesInMonth
	remMinutes = remMinutes % MinutesInMonth

	days := remMinutes / MinutesInDay
	finalMinutes := remMinutes % MinutesInDay

	calculatedDate := startDate.AddDate(int(years), int(months), int(days))
	return calculatedDate.Add(time.Minute * time.Duration(finalMinutes))
}
