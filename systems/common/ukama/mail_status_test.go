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

func TestMailStatus(t *testing.T) {
	t.Run("MailStatusValidString", func(tt *testing.T) {
		successMailStatus := ukama.ParseMailStatus("success")

		assert.NotNil(t, successMailStatus)
		assert.Equal(t, successMailStatus.String(), ukama.MailStatusSuccess.String())
		assert.Equal(t, uint8(successMailStatus), uint8(1))
	})

	t.Run("MailStatusValidNumber", func(tt *testing.T) {
		failedMailStatus := ukama.ParseMailStatus("2")

		assert.NotNil(t, failedMailStatus)
		assert.Equal(t, uint8(failedMailStatus), uint8(2))
		assert.Equal(t, failedMailStatus.String(), ukama.MailStatusFailed.String())
	})

	t.Run("MailStatusNonValidString", func(tt *testing.T) {
		unsupportedMailStatus := ukama.ParseMailStatus("failure")

		assert.NotNil(t, unsupportedMailStatus)
		assert.Equal(t, unsupportedMailStatus.String(), ukama.MailStatusPending.String())
		assert.Equal(t, uint8(unsupportedMailStatus), uint8(0))
	})

	t.Run("MailStatusNonValidNumber", func(tt *testing.T) {
		unknownMailStatus := ukama.MailStatus(uint8(10))

		assert.NotNil(t, unknownMailStatus)
		assert.Equal(t, unknownMailStatus.String(), ukama.MailStatusPending.String())
		assert.Equal(t, uint8(unknownMailStatus), uint8(10))
	})
}
