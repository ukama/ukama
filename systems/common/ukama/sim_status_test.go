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

func TestSimStatus(t *testing.T) {
	t.Run("SimTypeValidString", func(tt *testing.T) {
		simStatus := ukama.ParseSimStatus("active")

		assert.NotNil(t, simStatus)
		assert.Equal(t, simStatus.String(), ukama.SimStatusActive.String())
		assert.Equal(t, uint8(simStatus), uint8(1))
	})

	t.Run("SimTypeValidNumber", func(tt *testing.T) {
		simStatus := ukama.ParseSimStatus("3")

		assert.NotNil(t, simStatus)
		assert.Equal(t, uint8(simStatus), uint8(3))
		assert.Equal(t, simStatus.String(), ukama.SimStatusTerminated.String())
	})

	t.Run("SimTypeNonValidString", func(tt *testing.T) {
		simStatus := ukama.ParseSimStatus("failure")

		assert.NotNil(t, simStatus)
		assert.Equal(t, simStatus.String(), ukama.SimStatusUnknown.String())
		assert.Equal(t, uint8(simStatus), uint8(0))
	})

	t.Run("SimTypeNonValidNumber", func(tt *testing.T) {
		simStatus := ukama.SimStatus(uint8(10))

		assert.NotNil(t, simStatus)
		assert.Equal(t, simStatus.String(), ukama.SimStatusUnknown.String())
		assert.Equal(t, uint8(simStatus), uint8(10))
	})
}
