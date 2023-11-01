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
