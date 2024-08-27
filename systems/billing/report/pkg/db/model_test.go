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
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
)

func TestOwnerType(t *testing.T) {
	t.Run("OwnerTypeValidString", func(tt *testing.T) {
		owner := db.ParseOwnerType("org")

		assert.NotNil(t, owner)
		assert.Equal(t, owner.String(), "org")
		assert.Equal(t, uint8(owner), uint8(1))
	})

	t.Run("OwnerTypeValidNumber", func(tt *testing.T) {
		owner := db.ParseOwnerType("2")

		assert.NotNil(t, owner)
		assert.Equal(t, uint8(owner), uint8(2))
		assert.Equal(t, owner.String(), "subscriber")
	})

	t.Run("OwnerTypeNonValidString", func(tt *testing.T) {
		owner := db.ParseOwnerType("failure")

		assert.NotNil(t, owner)
		assert.Equal(t, owner.String(), "unknown")
		assert.Equal(t, uint8(owner), uint8(0))
	})

	t.Run("OwnerTypeNonValidNumber", func(tt *testing.T) {
		owner := db.OwnerType(uint8(10))

		assert.NotNil(t, owner)
		assert.Equal(t, owner.String(), "unknown")
		assert.Equal(t, uint8(owner), uint8(10))
	})
}

func TestReportType(t *testing.T) {
	t.Run("ReportTypeValidString", func(tt *testing.T) {
		report := db.ParseReportType("invoice")

		assert.NotNil(t, report)
		assert.Equal(t, report.String(), "invoice")
		assert.Equal(t, uint8(report), uint8(1))
	})

	t.Run("ReportTypeValidNumber", func(tt *testing.T) {
		report := db.ParseReportType("2")

		assert.NotNil(t, report)
		assert.Equal(t, uint8(report), uint8(2))
		assert.Equal(t, report.String(), "consumption")
	})

	t.Run("ReportTypeNonValidString", func(tt *testing.T) {
		report := db.ParseReportType("failure")

		assert.NotNil(t, report)
		assert.Equal(t, report.String(), "unknown")
		assert.Equal(t, uint8(report), uint8(0))
	})

	t.Run("ReportTypeNonValidNumber", func(tt *testing.T) {
		report := db.ReportType(uint8(10))

		assert.NotNil(t, report)
		assert.Equal(t, report.String(), "unknown")
		assert.Equal(t, uint8(report), uint8(10))
	})
}
