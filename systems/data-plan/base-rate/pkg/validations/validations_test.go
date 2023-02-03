package validations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// IsFutureDate Success case
func TestRateService_IsFutureDate_Success(t *testing.T) {
	assert.True(t, IsFutureDate("2023-10-10T00:00:00Z"))
}

// IsFutureDate Error case
func TestRateService_IsFutureDate_Error(t *testing.T) {
	assert.False(t, IsFutureDate("2023-10-10"))
	assert.False(t, IsFutureDate("2020-10-10T00:00:00Z"))
}

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
