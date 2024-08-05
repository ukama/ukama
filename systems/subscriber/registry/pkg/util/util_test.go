package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDOB(t *testing.T) {
	t.Run("Valid RFC3339 date", func(t *testing.T) {
		validDOB := "2024-08-01T12:34:56Z" // RFC3339 format
		expected := validDOB

		result, err := ValidateDOB(validDOB)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Invalid date format", func(t *testing.T) {
		invalidDOB := "2024/08/01 12:34:56" // Not RFC3339 format

		result, err := ValidateDOB(invalidDOB)

		assert.Error(t, err)
		assert.Equal(t, "invalid date format, must be RFC3339 standard", err.Error())
		assert.Empty(t, result)
	})

	t.Run("Empty date string", func(t *testing.T) {
		emptyDOB := ""

		result, err := ValidateDOB(emptyDOB)

		assert.Error(t, err)
		assert.Equal(t, "invalid date format, must be RFC3339 standard", err.Error())
		assert.Empty(t, result)
	})

	t.Run("Date with timezone offset", func(t *testing.T) {
		dobWithOffset := "2024-08-01T12:34:56+02:00" // RFC3339 with offset
		expected := "2024-08-01T12:34:56+02:00"

		result, err := ValidateDOB(dobWithOffset)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
