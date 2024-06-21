//go:build integration
// +build integration

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package integration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/config"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

const AppsPath = "/v1/hub/app"
const AppName = "hub-integration-test"
const AppVer = "0.1.2"
const FileName = "0.1.2.tar.gz"

type TestConfig struct {
	config.BaseConfig
	HubApigateway string
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{
		HubApigateway: "http://gateway:8080",
	}

	config.LoadConfig("integration", tConfig)
	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("%+v", tConfig)
}

// Call webhost endpoint and check response
func Test_AddCApp(t *testing.T) {

	appUrl := fmt.Sprintf("%s%s/%s/%s", tConfig.HubApigateway, AppsPath, AppName, AppVer)
	con := getFileContent(t)
	rest := resty.New().EnableTrace().SetDebug(tConfig.DebugMode)

	t.Run("Ping", func(t *testing.T) {
		r, err := rest.R().Get(tConfig.HubApigateway + "/ping")

		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode(), http.StatusOK)
	})

	t.Run("Put", func(tt *testing.T) {
		r, err := rest.R().SetHeader("Content-Type", "application/octet-stream").SetBody(con).Put(appUrl)

		assert.NoError(tt, err)
		log.Infof("Response: '%s'", r.String())
		assert.Equal(tt, http.StatusCreated, r.StatusCode())
	})

	t.Run("GetTarGz", func(tt *testing.T) {
		r, err := rest.R().Get(appUrl + ".tar.gz")
		if err != nil {
			assert.FailNow(tt, err.Error())
		}
		assert.Equal(tt, r.StatusCode(), http.StatusOK)
		assert.NoError(tt, err)

		body := r.Body()
		assert.Equal(tt, len(con), len(body))

		if !bytes.Equal(con, body) {
			assert.Fail(tt, "Expected file content is not equal to actual content")
		}
	})

	t.Run("GetChunkIndex", func(tt *testing.T) {
		for i := 0; i < 3; i++ {
			// wait for chunk to be created for 3 mins
			time.Sleep(time.Second * 60)
			log.Infof("Getting chunk index attempt %d", i)

			r, err := rest.R().Get(appUrl + ".caibx")
			if err != nil {
				assert.FailNow(tt, err.Error())
			}

			if r.StatusCode() == http.StatusNotFound {
				continue
			}

			assert.Equal(tt, r.StatusCode(), http.StatusOK)
			assert.NoError(tt, err)
		}
	})
}

func getFileContent(t *testing.T) []byte {
	f, err := os.Open("testdata/capp.tar.gz")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer f.Close()

	con, err := ioutil.ReadAll(f)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return con
}
