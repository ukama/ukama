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
	"testing"
	"time"
)

var (
	malformattedDate  = "2021/03/08"
	effectiveAtPast   = "2022-12-01T00:00:00Z"
	effectiveAtFuture = time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC3339)
)

func TestIsFutureDate(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		expected error
	}{
		{
			name:     "future date",
			date:     effectiveAtFuture,
			expected: nil,
		},
		{
			name:     "past date",
			date:     effectiveAtPast,
			expected: errors.New("Date is not in the future"),
		},
		{
			name:     "invalid date format",
			date:     malformattedDate,
			expected: errors.New("invalid date format, must be RFC3339 standard"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := IsFutureDate(tc.date)
			if err != nil && tc.expected == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && tc.expected != nil {
				t.Fatalf("expected error: %v, but got nil", tc.expected)
			}
		})
	}
}

func TestIsAfterDate(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		after    string
		expected error
	}{
		{
			name:     "invalid date format",
			date:     malformattedDate,
			after:    malformattedDate,
			expected: errors.New("invalid date format, must be RFC3339 standard"),
		},
		{
			name:     "invalid after format",
			date:     effectiveAtPast,
			after:    malformattedDate,
			expected: errors.New("invalid date format, must be RFC3339 standard"),
		},
		{
			name:     "date is after",
			date:     effectiveAtFuture,
			after:    effectiveAtPast,
			expected: nil,
		},
		{
			name:     "date is before",
			date:     effectiveAtPast,
			after:    effectiveAtFuture,
			expected: errors.New("Date is not in the future"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := IsAfterDate(tc.date, tc.after)
			if err != nil && tc.expected == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && tc.expected != nil {
				t.Fatalf("expected error: %v, but got nil", tc.expected)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	testCases := []struct {
		name     string
		date     string
		expected string
		hasError bool
	}{
		{
			name:     "valid date",
			date:     effectiveAtFuture,
			expected: effectiveAtFuture,
			hasError: false,
		},
		{
			name:     "invalid date format",
			date:     malformattedDate,
			expected: "",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValidateDate(tc.date)
			if tc.hasError && err == nil {
				t.Fatalf("expected error, but got nil")
			}
			if !tc.hasError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Fatalf("expected: %v, but got: %v", tc.expected, result)
			}
		})
	}
}
