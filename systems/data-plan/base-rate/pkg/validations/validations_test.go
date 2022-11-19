package validationErrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateService_IsFutureDate(t *testing.T) {
	//Success case
	assert.True(t, IsFutureDate("2023-10-10T00:00:00Z"))
	//Error case
	assert.False(t, IsFutureDate("2023-10-10"))
	assert.False(t, IsFutureDate("2020-10-10T00:00:00Z"))
}

func TestRateService_IsValidUploadReqArgs(t *testing.T) {
	//Success case
	assert.True(t, IsValidUploadReqArgs("https://test.com", "2023-10-10", "inter_mno_data"))
	//Error case
	assert.False(t, IsValidUploadReqArgs("", "2023-10-10", "inter_mno_data"))
	assert.False(t, IsValidUploadReqArgs("", "", "inter_mno_data"))
	assert.False(t, IsValidUploadReqArgs("", "", ""))
}
