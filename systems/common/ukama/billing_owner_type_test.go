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

func TestOwnerType(t *testing.T) {
	t.Run("OwnerTypeValidString", func(tt *testing.T) {
		owner := ukama.ParseOwnerType("org")

		assert.NotNil(t, owner)
		assert.Equal(t, owner.String(), ukama.OwnerTypeOrg.String())
		assert.Equal(t, uint8(owner), uint8(1))
	})

	t.Run("OwnerTypeValidNumber", func(tt *testing.T) {
		owner := ukama.ParseOwnerType("2")

		assert.NotNil(t, owner)
		assert.Equal(t, uint8(owner), uint8(2))
		assert.Equal(t, owner.String(), ukama.OwnerTypeSubscriber.String())
	})

	t.Run("OwnerTypeNonValidString", func(tt *testing.T) {
		owner := ukama.ParseOwnerType("failure")

		assert.NotNil(t, owner)
		assert.Equal(t, owner.String(), ukama.OwnerTypeUnknown.String())
		assert.Equal(t, uint8(owner), uint8(0))
	})

	t.Run("OwnerTypeNonValidNumber", func(tt *testing.T) {
		owner := ukama.OwnerType(uint8(10))

		assert.NotNil(t, owner)
		assert.Equal(t, owner.String(), ukama.OwnerTypeUnknown.String())
		assert.Equal(t, uint8(owner), uint8(10))
	})
}
