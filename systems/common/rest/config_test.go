/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"testing"

	"github.com/num30/config"
	"github.com/stretchr/testify/assert"
)

type testConf struct {
	Conf HttpConfig
}

type testConfPoint struct {
	Conf *HttpConfig `default:"{}"`
}

func Test_DefaultValues(t *testing.T) {

	testConf := &testConf{}
	reader := config.NewConfReader("test")
	err := reader.Read(testConf)
	if assert.NoError(t, err) {
		assert.Equal(t, "http://localhost", testConf.Conf.Cors.AllowOrigins[0])
	}
}

func Test_DefaultValuesPointer(t *testing.T) {

	testConf := &testConfPoint{}
	reader := config.NewConfReader("test")
	err := reader.Read(testConf)
	if assert.NoError(t, err) {
		assert.Equal(t, "http://localhost", testConf.Conf.Cors.AllowOrigins[0])
	}
}
