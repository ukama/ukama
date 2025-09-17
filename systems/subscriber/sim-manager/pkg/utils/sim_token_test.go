/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

const (
	testIccid    = "8910300000003540781"
	textToEncode = `{"iccid": ` + testIccid + `}`
	testKey1     = "the-key-has-to-be-32-bytes-long!"
	testKey2     = "that--key-is-also-32-bytes-long!"
)

func Test_IccidTokenConversion(t *testing.T) {
	t.Run("IccidRetrieved", func(t *testing.T) {
		token, err := GenerateTokenFromIccid(testIccid, testKey2)
		assert.NoError(t, err)
		res, err := GetIccidFromToken(token, testKey2)
		assert.NoError(t, err)
		assert.Equal(t, testIccid, res)
	})

	t.Run("IccidNotRetrieved", func(t *testing.T) {
		token, err := GenerateTokenFromIccid(testIccid, testKey1)
		assert.NoError(t, err)

		res, err := GetIccidFromToken(token, testKey2)
		assert.Error(t, err)
		assert.NotEqual(t, testIccid, res)
	})
}

func Test_encrypt(t *testing.T) {
	r, err := encrypt(textToEncode, testKey1)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "encrypt failed")
	}

	log.Infof("encrypted sim: %s", r)

	res, err := decrypt(r, testKey1)
	assert.NoError(t, err)
	assert.Equal(t, textToEncode, res)
}

func Test_encryptErrors(t *testing.T) {
	tests := []struct {
		key  string
		text string
	}{
		{key: "short testKey", text: textToEncode},
		{key: "", text: textToEncode},
	}
	for _, tt := range tests {
		_, err := encrypt(tt.text, tt.key)
		assert.Error(t, err)
	}
}

func Test_decryptErrors(t *testing.T) {
	encrText, err := encrypt(textToEncode, testKey1)
	if err != nil {
		assert.FailNow(t, "encrypt failed", err)
	}

	tests := []struct {
		key  string
		text string
	}{
		{key: "short key", text: encrText},
		{key: "", text: textToEncode},
		{key: "", text: "Not encoded text"},
	}
	for _, tt := range tests {
		_, err := decrypt(tt.text, tt.key)
		assert.Error(t, err)
	}
}
