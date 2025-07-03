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

func TestMessageUnitType(t *testing.T) {
	t.Run("MessageUnitTypeValidString", func(tt *testing.T) {
		intMessageUnitType := ukama.ParseMessageType("int")

		assert.NotNil(t, intMessageUnitType)
		assert.Equal(t, intMessageUnitType.String(), ukama.MessageUnitTypeInt.String())
		assert.Equal(t, uint8(intMessageUnitType), uint8(1))
		assert.Equal(t, ukama.ReturnMessageUnits(intMessageUnitType), int64(1))
	})

	t.Run("MessageUnitTypeValidNumber", func(tt *testing.T) {
		intMessageUnitType := ukama.ParseMessageType("1")

		assert.NotNil(t, intMessageUnitType)
		assert.Equal(t, uint8(intMessageUnitType), uint8(1))
		assert.Equal(t, intMessageUnitType.String(), ukama.MessageUnitTypeInt.String())
		assert.Equal(t, ukama.ReturnMessageUnits(intMessageUnitType), int64(1))
	})

	t.Run("MessageUnitTypeNonValidString", func(tt *testing.T) {
		unsupportedMessageUnitType := ukama.ParseMessageType("failure")

		assert.NotNil(t, unsupportedMessageUnitType)
		assert.Equal(t, unsupportedMessageUnitType.String(), ukama.MessageUnitTypeUnknown.String())
		assert.Equal(t, uint8(unsupportedMessageUnitType), uint8(0))
		assert.Equal(t, ukama.ReturnMessageUnits(unsupportedMessageUnitType), int64(0))
	})

	t.Run("MessageUnitTypeNonValidNumber", func(tt *testing.T) {
		unknownMessageUnitType := ukama.MessageUnitType(uint8(10))

		assert.NotNil(t, unknownMessageUnitType)
		assert.Equal(t, unknownMessageUnitType.String(), ukama.MessageUnitTypeUnknown.String())
		assert.Equal(t, uint8(unknownMessageUnitType), uint8(10))
		assert.Equal(t, ukama.ReturnMessageUnits(unknownMessageUnitType), int64(0))
	})
}
