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
	var mockeEffectiveAt = time.Now().Add(time.Hour * 24 * 7 * time.Duration(4)).Format(time.RFC3339Nano)

    testCases := []struct {
        name     string
        date     string
        expected error
    }{
        {
            name:     mockeEffectiveAt,
            date:     "Mon, 08 Mar 2023 00:00:00 GMT",
            expected: nil,
        },
        {
            name:     "past date",
            date:     "Mon, 01 Mar 2021 00:00:00 GMT",
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
    testCases := []struct {
        name     string
        date     string
        expected string
        hasError bool
    }{
        {
            name:     "valid date",
            date:     "08-03-2023",
            expected: "Wed, 08 Mar 2023 00:00:00 UTC",
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