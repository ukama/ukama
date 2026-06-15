/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"
)

const (
	testIccid    = "89860012345678901000"
	textToEncode = `{"iccid": ` + testIccid + `}`
	testKey1     = "the-key-has-to-be-32-bytes-long!"
	testKey2     = "that--key-is-also-32-bytes-long!"
)

func Test_IccidTokenConversion(t *testing.T) {
	tokenCodec1 := utils.NewTokenCodec(testKey1)
	tokenCodec2 := utils.NewTokenCodec(testKey2)

	t.Run("IccidRetrieved", func(t *testing.T) {
		token, err := tokenCodec1.GenerateTokenFromIccid(testIccid)
		assert.NoError(t, err)
		res, err := tokenCodec1.GetIccidFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, testIccid, res)
	})

	t.Run("IccidNotRetrieved", func(t *testing.T) {
		token, err := tokenCodec1.GenerateTokenFromIccid(testIccid)
		assert.NoError(t, err)

		res, err := tokenCodec2.GetIccidFromToken(token)
		assert.Error(t, err)
		assert.NotEqual(t, testIccid, res)
	})
}
