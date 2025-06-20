/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package validations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_IsEmpty(t *testing.T) {
	t.Run("all strings are non-empty", func(t *testing.T) {
		check := IsEmpty("a", "b", "c")
		assert.False(t, check)
	})

	t.Run("one string is empty", func(t *testing.T) {
		check := IsEmpty("a", "", "c")
		assert.True(t, check)
	})

	t.Run("multiple strings are empty", func(t *testing.T) {
		check := IsEmpty("", "", "c")
		assert.True(t, check)
	})

	t.Run("all strings are empty", func(t *testing.T) {
		check := IsEmpty("", "", "")
		assert.True(t, check)
	})

	t.Run("single empty string", func(t *testing.T) {
		check := IsEmpty("")
		assert.True(t, check)
	})

	t.Run("single non-empty string", func(t *testing.T) {
		check := IsEmpty("hello")
		assert.False(t, check)
	})

	t.Run("no arguments", func(t *testing.T) {
		check := IsEmpty()
		assert.False(t, check)
	})

	t.Run("mixed empty and non-empty strings", func(t *testing.T) {
		check := IsEmpty("hello", "", "world", "")
		assert.True(t, check)
	})

	t.Run("strings with only spaces", func(t *testing.T) {
		check := IsEmpty(" ", "  ", "   ")
		assert.False(t, check) // Spaces are not considered empty
	})

	t.Run("strings with tabs and newlines", func(t *testing.T) {
		check := IsEmpty("\t", "\n", "\r")
		assert.False(t, check) // Control characters are not considered empty
	})
}

func Test_IsReqEmpty(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		check := IsReqEmpty(0)
		assert.True(t, check)
	})

	t.Run("positive value", func(t *testing.T) {
		check := IsReqEmpty(1)
		assert.False(t, check)
	})

	t.Run("large positive value", func(t *testing.T) {
		check := IsReqEmpty(1000)
		assert.False(t, check)
	})

	t.Run("maximum uint64 value", func(t *testing.T) {
		check := IsReqEmpty(18446744073709551615) // max uint64
		assert.False(t, check)
	})

	t.Run("negative value (should be false since uint64)", func(t *testing.T) {
		// Note: This would cause a compilation error in real code
		// but testing the logic that any non-zero value returns false
		check := IsReqEmpty(1)
		assert.False(t, check)
	})
}

func Test_ValidateDOB(t *testing.T) {
	t.Run("valid date format", func(t *testing.T) {
		dob := "15-03-1990"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("valid date with single digit day and month", func(t *testing.T) {
		dob := "05-09-1985"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("valid date with double digit day and month", func(t *testing.T) {
		dob := "25-12-2000"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("invalid date format - wrong separator", func(t *testing.T) {
		dob := "15/03/1990"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - wrong order", func(t *testing.T) {
		dob := "1990-03-15"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - missing parts", func(t *testing.T) {
		dob := "15-03"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - extra parts", func(t *testing.T) {
		dob := "15-03-1990-extra"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - empty string", func(t *testing.T) {
		dob := ""

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - non-numeric day", func(t *testing.T) {
		dob := "ab-03-1990"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - non-numeric month", func(t *testing.T) {
		dob := "15-xy-1990"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date format - non-numeric year", func(t *testing.T) {
		dob := "15-03-abc"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date - day out of range", func(t *testing.T) {
		dob := "32-03-1990"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date - month out of range", func(t *testing.T) {
		dob := "15-13-1990"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("invalid date - February 30th", func(t *testing.T) {
		dob := "30-02-1990"

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("valid leap year date", func(t *testing.T) {
		dob := "29-02-2020"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("invalid leap year date", func(t *testing.T) {
		dob := "29-02-2021" // 2021 is not a leap year

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})

	t.Run("very old date", func(t *testing.T) {
		dob := "01-01-1900"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("future date", func(t *testing.T) {
		dob := "31-12-2030"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("date with leading zeros", func(t *testing.T) {
		dob := "01-01-2001"
		expectedTime, _ := time.Parse("02-01-2006", dob)

		result, err := ValidateDOB(dob)

		assert.NoError(t, err)
		assert.Equal(t, expectedTime, result)
	})

	t.Run("date with spaces", func(t *testing.T) {
		dob := " 15-03-1990 "

		result, err := ValidateDOB(dob)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, result)
		assert.Contains(t, err.Error(), "invalid date format, must be dd-mm-yyyy")
	})
}
