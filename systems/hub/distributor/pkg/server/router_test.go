/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ukama/ukama/systems/hub/distributor/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	log "github.com/sirupsen/logrus"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
)

func init() {
	pkg.IsDebugMode = true
}

func Test_RouterPing(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	defconf := pkg.NewConfig(pkg.ServiceName)

	r := NewRouter(defconf, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func Test_RouterPut(t *testing.T) {
	// arrange
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	w := httptest.NewRecorder()
	f := []byte(`{ "store":"./test/data/art" }`)

	req, _ := http.NewRequest("PUT", ChunksPath+"/ukamaos/1.0.1", bytes.NewBuffer(f))

	defconf := pkg.NewConfig(pkg.ServiceName)
	defconf.Distribution.Chunk.Stores[0] = "./test/data/store"

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
	r := NewRouter(defconf, msgbusClient).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 200, w.Code)

	putIndexFileContent(t, w.Body)
}

func Test_RouterPutNoStore(t *testing.T) {
	// arrange
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", ChunksPath+"/ukamaos/1.0.1", nil)

	defconf := pkg.NewConfig(pkg.ServiceName)

	r := NewRouter(defconf, nil).fizz.Engine()

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Error:Field validation for 'Store'")
}

func putIndexFileContent(tt *testing.T, br io.Reader) {
	tt.Helper()

	file := "./test/data/index/index.caidx"

	f, err := os.Create(file)
	if err != nil {
		assert.FailNow(tt, err.Error())
	}
	defer f.Close()

	bytes, err := io.Copy(f, br)
	if err != nil {
		assert.FailNow(tt, err.Error())
	}

	if bytes <= 0 {
		assert.FailNow(tt, "expected file contents but looks like its empty")
	}

	log.Debugf("Index file %s created with %d bytes.", file, bytes)
}
