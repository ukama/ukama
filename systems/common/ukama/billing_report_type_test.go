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

func TestReportType(t *testing.T) {
	t.Run("ReportTypeValidString", func(tt *testing.T) {
		report := ukama.ParseReportType("invoice")

		assert.NotNil(t, report)
		assert.Equal(t, report.String(), ukama.ReportTypeInvoice.String())
		assert.Equal(t, uint8(report), uint8(1))
	})

	t.Run("ReportTypeValidNumber", func(tt *testing.T) {
		report := ukama.ParseReportType("2")

		assert.NotNil(t, report)
		assert.Equal(t, uint8(report), uint8(2))
		assert.Equal(t, report.String(), ukama.ReportTypeConsumption.String())
	})

	t.Run("ReportTypeNonValidString", func(tt *testing.T) {
		report := ukama.ParseReportType("failure")

		assert.NotNil(t, report)
		assert.Equal(t, report.String(), ukama.ReportTypeUnknown.String())
		assert.Equal(t, uint8(report), uint8(0))
	})

	t.Run("ReportTypeNonValidNumber", func(tt *testing.T) {
		report := ukama.ReportType(uint8(10))

		assert.NotNil(t, report)
		assert.Equal(t, report.String(), ukama.ReportTypeUnknown.String())
		assert.Equal(t, uint8(report), uint8(10))
	})
}
