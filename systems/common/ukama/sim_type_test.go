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

func TestSimType(t *testing.T) {
	t.Run("SimTypeValidString", func(tt *testing.T) {
		simType := ukama.ParseSimType("test")

		assert.NotNil(t, simType)
		assert.Equal(t, simType.String(), "test")
		assert.Equal(t, uint8(simType), uint8(1))
	})

	t.Run("SimTypeValidNumber", func(tt *testing.T) {
		simType := ukama.ParseSimType("3")

		assert.NotNil(t, simType)
		assert.Equal(t, uint8(simType), uint8(3))
		assert.Equal(t, simType.String(), "ukama_data")
	})

	t.Run("SimTypeNonValidString", func(tt *testing.T) {
		simType := ukama.ParseSimType("failure")

		assert.NotNil(t, simType)
		assert.Equal(t, simType.String(), "unknown")
		assert.Equal(t, uint8(simType), uint8(0))
	})

	t.Run("SimTypeNonValidNumber", func(tt *testing.T) {
		simType := ukama.SimType(uint8(10))

		assert.NotNil(t, simType)
		assert.Equal(t, simType.String(), "unknown")
		assert.Equal(t, uint8(simType), uint8(10))
	})
}

func TestSimStatus(t *testing.T) {
	t.Run("SimTypeValidString", func(tt *testing.T) {
		simStatus := ukama.ParseSimStatus("active")

		assert.NotNil(t, simStatus)
		assert.Equal(t, simStatus.String(), "active")
		assert.Equal(t, uint8(simStatus), uint8(1))
	})

	t.Run("SimTypeValidNumber", func(tt *testing.T) {
		simStatus := ukama.ParseSimStatus("3")

		assert.NotNil(t, simStatus)
		assert.Equal(t, uint8(simStatus), uint8(3))
		assert.Equal(t, simStatus.String(), "terminated")
	})

	t.Run("SimTypeNonValidString", func(tt *testing.T) {
		simStatus := ukama.ParseSimStatus("failure")

		assert.NotNil(t, simStatus)
		assert.Equal(t, simStatus.String(), "unknown")
		assert.Equal(t, uint8(simStatus), uint8(0))
	})

	t.Run("SimTypeNonValidNumber", func(tt *testing.T) {
		simStatus := ukama.SimStatus(uint8(10))

		assert.NotNil(t, simStatus)
		assert.Equal(t, simStatus.String(), "unknown")
		assert.Equal(t, uint8(simStatus), uint8(10))
	})
}
