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

func TestDataUnitType(t *testing.T) {
	t.Run("DataUnitTypeValidString", func(tt *testing.T) {
		dataUnitType := ukama.ParseDataUnitType("Bytes")

		assert.NotNil(t, dataUnitType)
		assert.Equal(t, dataUnitType.String(), ukama.DataUnitTypeB.String())
		assert.Equal(t, uint8(dataUnitType), uint8(1))
		assert.Equal(t, ukama.ReturnDataUnits(dataUnitType), int64(1))
	})

	t.Run("DataUnitTypeValidNumber", func(tt *testing.T) {
		dataUnitType := ukama.ParseDataUnitType("2")

		assert.NotNil(t, dataUnitType)
		assert.Equal(t, uint8(dataUnitType), uint8(2))
		assert.Equal(t, dataUnitType.String(), ukama.DataUnitTypeKB.String())
		assert.Equal(t, ukama.ReturnDataUnits(dataUnitType), int64(1))
	})

	t.Run("DataUnitTypeNonValidString", func(tt *testing.T) {
		dataUnitType := ukama.ParseDataUnitType("failure")

		assert.NotNil(t, dataUnitType)
		assert.Equal(t, dataUnitType.String(), ukama.DataUnitTypeUnknown.String())
		assert.Equal(t, uint8(dataUnitType), uint8(0))
		assert.Equal(t, ukama.ReturnDataUnits(dataUnitType), int64(0))
	})

	t.Run("DataUnitTypeNonValidNumber", func(tt *testing.T) {
		dataUnitType := ukama.DataUnitType(uint8(10))

		assert.NotNil(t, dataUnitType)
		assert.Equal(t, dataUnitType.String(), ukama.DataUnitTypeUnknown.String())
		assert.Equal(t, uint8(dataUnitType), uint8(10))
		assert.Equal(t, ukama.ReturnDataUnits(dataUnitType), int64(0))
	})
}
