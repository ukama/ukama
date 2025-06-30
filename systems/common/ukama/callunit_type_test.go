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

func TestCallUnitType(t *testing.T) {
	t.Run("CallUnitTypeValidString", func(tt *testing.T) {
		callUnitType := ukama.ParseCallUnitType("seconds")

		assert.NotNil(t, callUnitType)
		assert.Equal(t, callUnitType.String(), ukama.CallUnitTypeSec.String())
		assert.Equal(t, uint8(callUnitType), uint8(1))
		assert.Equal(t, ukama.ReturnCallUnits(callUnitType), int64(1))
	})

	t.Run("CallUnitTypeValidNumber", func(tt *testing.T) {
		callUnitType := ukama.ParseCallUnitType("2")

		assert.NotNil(t, callUnitType)
		assert.Equal(t, uint8(callUnitType), uint8(2))
		assert.Equal(t, callUnitType.String(), ukama.CallUnitTypeMin.String())
		assert.Equal(t, ukama.ReturnCallUnits(callUnitType), int64(60))
	})

	t.Run("CallUnitTypeNonValidString", func(tt *testing.T) {
		callUnitType := ukama.ParseCallUnitType("failure")

		assert.NotNil(t, callUnitType)
		assert.Equal(t, callUnitType.String(), ukama.CallUnitTypeUnknown.String())
		assert.Equal(t, uint8(callUnitType), uint8(0))
		assert.Equal(t, ukama.ReturnCallUnits(callUnitType), int64(0))
	})

	t.Run("CallUnitTypeNonValidNumber", func(tt *testing.T) {
		callUnitType := ukama.CallUnitType(uint8(10))

		assert.NotNil(t, callUnitType)
		assert.Equal(t, callUnitType.String(), ukama.CallUnitTypeUnknown.String())
		assert.Equal(t, uint8(callUnitType), uint8(10))
		assert.Equal(t, ukama.ReturnCallUnits(callUnitType), int64(0))
	})
}
