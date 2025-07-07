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

func TestComponentCategory(t *testing.T) {
	t.Run("ComponentCategoryValidString", func(tt *testing.T) {
		//TODO: add case insensitive test case after ParseType function updates.
		accessCategory := ukama.ParseType("access")

		assert.NotNil(t, accessCategory)
		assert.Equal(t, accessCategory.String(), ukama.ACCESS.String())
		assert.Equal(t, uint8(accessCategory), uint8(1))
	})

	t.Run("ComponentCategoryValidNumber", func(tt *testing.T) {
		backhaulCategory := ukama.ParseType("2")

		assert.NotNil(t, backhaulCategory)
		assert.Equal(t, uint8(backhaulCategory), uint8(2))
		assert.Equal(t, backhaulCategory.String(), ukama.BACKHAUL.String())
	})

	t.Run("ComponentCategoryNonValidString", func(tt *testing.T) {
		defaultCategory := ukama.ParseType("failure")

		assert.NotNil(t, defaultCategory)
		assert.Equal(t, defaultCategory.String(), ukama.ALL.String())
		assert.Equal(t, uint8(defaultCategory), uint8(0))
	})

	t.Run("ComponentCategoryNonValidNumber", func(tt *testing.T) {
		defaultCategory := ukama.ComponentCategory(uint8(10))

		assert.NotNil(t, defaultCategory)
		assert.Equal(t, defaultCategory.String(), ukama.ALL.String())
		assert.Equal(t, uint8(defaultCategory), uint8(10))
	})
}
