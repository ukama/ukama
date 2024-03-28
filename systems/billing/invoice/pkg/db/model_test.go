/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/db"
)

func TestInvoiceeType(t *testing.T) {
	t.Run("InvoiceeTypeValidString", func(tt *testing.T) {
		status := db.ParseInvoiceeType("org")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "org")
		assert.Equal(t, uint8(status), uint8(1))
	})

	t.Run("InvoiceeTypeValidNumber", func(tt *testing.T) {
		status := db.ParseInvoiceeType("2")

		assert.NotNil(t, status)
		assert.Equal(t, uint8(status), uint8(2))
		assert.Equal(t, status.String(), "subscriber")
	})

	t.Run("InvoiceeTypeNonValidString", func(tt *testing.T) {
		status := db.ParseInvoiceeType("failure")

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "unknown")
		assert.Equal(t, uint8(status), uint8(0))
	})

	t.Run("InvoiceeTypeNonValidNumber", func(tt *testing.T) {
		status := db.InvoiceeType(uint8(10))

		assert.NotNil(t, status)
		assert.Equal(t, status.String(), "unknown")
		assert.Equal(t, uint8(status), uint8(10))
	})
}
