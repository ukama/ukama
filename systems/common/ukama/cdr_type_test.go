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

func TestCdrType(t *testing.T) {
	t.Run("CdrTypeValidString", func(tt *testing.T) {
		cdr := ukama.ParseCdrType("datA")

		assert.NotNil(t, cdr)
		assert.Equal(t, cdr.String(), ukama.CdrTypeData.String())
		assert.Equal(t, uint8(cdr), uint8(1))
	})

	t.Run("CdrTypeValidNumber", func(tt *testing.T) {
		cdr := ukama.ParseCdrType("3")

		assert.NotNil(t, cdr)
		assert.Equal(t, uint8(cdr), uint8(3))
		assert.Equal(t, cdr.String(), ukama.CdrTypeSms.String())
	})

	t.Run("CdrTypeNonValidString", func(tt *testing.T) {
		cdr := ukama.ParseCdrType("failure")

		assert.NotNil(t, cdr)
		assert.Equal(t, cdr.String(), ukama.CdrTypeUnknown.String())
		assert.Equal(t, uint8(cdr), uint8(0))
	})

	t.Run("CdrTypeNonValidNumber", func(tt *testing.T) {
		cdr := ukama.CdrType(uint8(10))

		assert.NotNil(t, cdr)
		assert.Equal(t, cdr.String(), ukama.CdrTypeUnknown.String())
		assert.Equal(t, uint8(cdr), uint8(10))
	})
}
