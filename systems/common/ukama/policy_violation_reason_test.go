/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/ukama"
)

func TestPolicyViolationReason(t *testing.T) {
	t.Run("PolicyViolationReasonValidString", func(tt *testing.T) {
		r := ukama.ParsePolicyViolationReason("Data_Cap_Exceeded")

		assert.NotNil(t, r)
		assert.Equal(t, r.String(), ukama.PolicyViolationReasonDataCapExceeded.String())
		assert.Equal(t, uint8(r), uint8(1))
	})

	t.Run("PolicyViolationReasonValidNumber", func(tt *testing.T) {
		r := ukama.ParsePolicyViolationReason("2")

		assert.NotNil(t, r)
		assert.Equal(t, uint8(r), uint8(2))
		assert.Equal(t, r.String(), ukama.PolicyViolationReasonPackageExpired.String())
	})

	t.Run("PolicyViolationReasonNonValidString", func(tt *testing.T) {
		r := ukama.ParsePolicyViolationReason("failure")

		assert.NotNil(t, r)
		assert.Equal(t, r.String(), ukama.PolicyViolationReasonUnknown.String())
		assert.Equal(t, uint8(r), uint8(0))
	})

	t.Run("PolicyViolationReasonNonValidNumber", func(tt *testing.T) {
		r := ukama.PolicyViolationReason(uint8(10))

		assert.NotNil(t, r)
		assert.Equal(t, r.String(), ukama.PolicyViolationReasonUnknown.String())
		assert.Equal(t, uint8(r), uint8(10))
	})
}
