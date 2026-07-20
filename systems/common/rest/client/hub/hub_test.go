/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package hub_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/hub"
)

const baseURL = "http://test-hub-api-gw.com"

func respondWith(status int, body string, assertReq func(*http.Request)) client.RoundTripFunc {
	return func(req *http.Request) *http.Response {
		if assertReq != nil {
			assertReq(req)
		}
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		}
	}
}

func TestHubClient_ListApps(t *testing.T) {
	t.Run("Found", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusOK, `{"artifact":["app1","app2"]}`, func(req *http.Request) {
			assert.Equal(tt, baseURL+hub.HubEndpoint+"/app", req.URL.String())
		}))

		apps, err := c.ListApps("app")

		assert.NoError(tt, err)
		assert.Equal(tt, []string{"app1", "app2"}, apps)
	})

	t.Run("NotFoundReturnsEmpty", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusNotFound, `{"error":"not found"}`, nil))

		apps, err := c.ListApps("app")

		assert.NoError(tt, err)
		assert.Empty(tt, apps)
	})
}

func TestHubClient_ListVersions(t *testing.T) {
	// The Hub api-gateway serializes int64 (size) as a JSON *string* — this is the
	// case that previously broke unmarshalling.
	t.Run("StringEncodedSize", func(tt *testing.T) {
		body := `{"versions":[{"version":"1.1.1-manual","FormatInfo":[` +
			`{"type":"tar.gz","size":"1700"},{"type":"chunk","size":"1700"}]}]}`

		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusOK, body, func(req *http.Request) {
			assert.Equal(tt, baseURL+hub.HubEndpoint+"/app/example", req.URL.String())
		}))

		versions, err := c.ListVersions("example", "app")

		assert.NoError(tt, err)
		assert.Len(tt, versions, 1)
		assert.Equal(tt, "example", versions[0].Name)
		assert.Equal(tt, "app", versions[0].Type)
		assert.Equal(tt, "1.1.1-manual", versions[0].Version)
		assert.Equal(tt, int64(1700), versions[0].SizeBytes)
		assert.True(tt, versions[0].Chunked)
	})

	t.Run("NumericSizeAndNoChunk", func(tt *testing.T) {
		body := `{"versions":[{"version":"0.0.1","FormatInfo":[{"type":"tar.gz","size":42}]}]}`

		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusOK, body, nil))

		versions, err := c.ListVersions("example", "app")

		assert.NoError(tt, err)
		assert.Len(tt, versions, 1)
		assert.Equal(tt, int64(42), versions[0].SizeBytes)
		assert.False(tt, versions[0].Chunked)
	})

	t.Run("NotFoundReturnsNil", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusNotFound, `{"error":"Artifact name is not valid"}`, nil))

		versions, err := c.ListVersions("missing", "app")

		assert.NoError(tt, err)
		assert.Nil(tt, versions)
	})

	t.Run("ServerErrorReturnsError", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusInternalServerError, `{"error":"boom"}`, nil))

		versions, err := c.ListVersions("example", "app")

		assert.Error(tt, err)
		assert.Nil(tt, versions)
	})
}

func TestHubClient_VersionExists(t *testing.T) {
	body := `{"versions":[{"version":"1.0.0","FormatInfo":[{"type":"tar.gz","size":"10"}]},` +
		`{"version":"1.1.1-manual","FormatInfo":[{"type":"tar.gz","size":"20"}]}]}`

	t.Run("Present", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusOK, body, nil))

		ok, err := c.VersionExists("example", "app", "1.1.1-manual")

		assert.NoError(tt, err)
		assert.True(tt, ok)
	})

	t.Run("Absent", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusOK, body, nil))

		ok, err := c.VersionExists("example", "app", "9.9.9")

		assert.NoError(tt, err)
		assert.False(tt, ok)
	})

	t.Run("NotFoundIsFalseNoError", func(tt *testing.T) {
		c := hub.NewHubClient(baseURL)
		c.R.C.SetTransport(respondWith(http.StatusNotFound, `{"error":"not found"}`, nil))

		ok, err := c.VersionExists("missing", "app", "1.0.0")

		assert.NoError(tt, err)
		assert.False(tt, ok)
	})
}
