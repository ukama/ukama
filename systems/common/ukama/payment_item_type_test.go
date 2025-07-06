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

func TestItemType(t *testing.T) {
	t.Run("ItemTypeValidString", func(tt *testing.T) {
		packageType := ukama.ParseItemType("PacKagE")

		assert.NotNil(t, packageType)
		assert.Equal(t, packageType.String(), ukama.ItemTypePackage.String())
		assert.Equal(t, uint8(packageType), uint8(1))
	})

	t.Run("ItemTypeValidNumber", func(tt *testing.T) {
		invoiceType := ukama.ParseItemType("2")

		assert.NotNil(t, invoiceType)
		assert.Equal(t, uint8(invoiceType), uint8(2))
		assert.Equal(t, invoiceType.String(), ukama.ItemTypeInvoice.String())
	})

	t.Run("ItemTypeNonValidString", func(tt *testing.T) {
		unsupportedType := ukama.ParseItemType("failure")

		assert.NotNil(t, unsupportedType)
		assert.Equal(t, unsupportedType.String(), ukama.ItemTypeUnknown.String())
		assert.Equal(t, uint8(unsupportedType), uint8(0))
	})

	t.Run("ItemTypeNonValidNumber", func(tt *testing.T) {
		unknownType := ukama.ItemType(uint8(10))

		assert.NotNil(t, unknownType)
		assert.Equal(t, unknownType.String(), ukama.ItemTypeUnknown.String())
		assert.Equal(t, uint8(unknownType), uint8(10))
	})
}
