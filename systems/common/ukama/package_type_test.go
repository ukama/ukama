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

func TestPackageType(t *testing.T) {
	t.Run("PackageTypeValidString", func(tt *testing.T) {
		prepaidType := ukama.ParsePackageType("prepaiD")

		assert.NotNil(t, prepaidType)
		assert.Equal(t, prepaidType.String(), "prepaid")
		assert.Equal(t, uint8(prepaidType), uint8(1))
	})

	t.Run("PackageTypeValidNumber", func(tt *testing.T) {
		postpaidType := ukama.ParsePackageType("2")

		assert.NotNil(t, postpaidType)
		assert.Equal(t, uint8(postpaidType), uint8(2))
		assert.Equal(t, postpaidType.String(), "postpaid")
	})

	t.Run("PackageTypeNonValidString", func(tt *testing.T) {
		unsupportedType := ukama.ParsePackageType("failure")

		assert.NotNil(t, unsupportedType)
		assert.Equal(t, unsupportedType.String(), "unknown")
		assert.Equal(t, uint8(unsupportedType), uint8(0))
	})

	t.Run("PackageTypeNonValidNumber", func(tt *testing.T) {
		unknownType := ukama.PackageType(uint8(10))

		assert.NotNil(t, unknownType)
		assert.Equal(t, unknownType.String(), "unknown")
		assert.Equal(t, uint8(unknownType), uint8(10))
	})
}
