package validations

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// IsValidUploadReqArgs Success case
func TestRateService_IsValidUploadReqArgs_Success(t *testing.T) {
	assert.True(t, IsValidUploadReqArgs("https://test.com", "2023-10-10", "INTER_MNO_DATA"))
}

// IsValidUploadReqArgs Error case
func TestRateService_IsValidUploadReqArgs_Error(t *testing.T) {
	assert.False(t, IsValidUploadReqArgs("", "2023-10-10", "INTER_MNO_DATA"))
	assert.False(t, IsValidUploadReqArgs("", "", "INTER_MNO_DATA"))
	assert.False(t, IsValidUploadReqArgs("", "", ""))
}
func TestIsFutureDate(t *testing.T) {
	var mockeEffectiveAtFuture = time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC1123)
	var mockeEffectiveAtPast = time.Now().Add(time.Hour * 20 * 365 * 15).Format(time.RFC3339)

	testCases := []struct {
		name     string
		date     string
		expected error
	}{
		{
			name:     "future date",
			date:     mockeEffectiveAtFuture,
			expected: nil,
		},
		{
			name:     "past date",
			date:     mockeEffectiveAtPast,
			expected: errors.New("Date is not in the future"),
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

func TestValidateDate(t *testing.T) {
	var mockeEffectiveAtFuture = time.Now().Add(time.Hour * 24 * 365 * 15).Format(time.RFC1123)
	testCases := []struct {
		name     string
		date     string
		expected string
		hasError bool
	}{
		{
			name:     "valid date",
			date:     mockeEffectiveAtFuture,
			expected: mockeEffectiveAtFuture,
			hasError: false,
		},
		{
			name:     "invalid date format",
			date:     "2021/03/08",
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
