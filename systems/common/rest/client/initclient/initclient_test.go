/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package initclient_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/initclient"
	"github.com/ukama/ukama/systems/common/uuid"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
)

const (
	orgName    = "test-org"
	systemName = "test-system"
	testIp     = "10.0.0.1"
	testPort   = 9090
)

func TestInitClient_GetSystem(t *testing.T) {
	t.Run("InitFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Equal(tt, req.URL.String(), initclient.InitApiEndpoint+"/"+orgName+"/systems/"+systemName)

			// fake systemIP info
			systemIP := `{"systemId": "03cb753f-5e03-4c97-8e47-625115476c72", "orgName": "test-org", "systemName": "test-system"}`

			// Send mock response
			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",

				// Send response to be tested
				Body: io.NopCloser(bytes.NewBufferString(systemIP)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		}

		testInitClient := initclient.NewInitClient("")

		// We replace the transport mechanism by mocking the http request
		// so that the test stays a unit test e.g, no server/network call.
		testInitClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testInitClient.GetSystem(orgName, systemName)

		assert.NoError(tt, err)
		assert.Equal(tt, orgName, s.OrgName)
	})

	t.Run("InitNotFound", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), initclient.InitApiEndpoint+"/"+orgName+"/systems/"+systemName)

			// error payload
			resp := `{"error":"not found"}`

			return &http.Response{
				StatusCode: 404,
				Status:     "404 NOT FOUND",
				Body:       io.NopCloser(bytes.NewBufferString(resp)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}
		}

		testInitClient := initclient.NewInitClient("")

		testInitClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testInitClient.GetSystem(orgName, systemName)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidResponsePayload", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), initclient.InitApiEndpoint+"/"+orgName+"/systems/"+systemName)

			return &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		testInitClient := initclient.NewInitClient("")

		testInitClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testInitClient.GetSystem(orgName, systemName)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("RequestFailure", func(tt *testing.T) {
		mockTransport := func(req *http.Request) *http.Response {
			assert.Equal(tt, req.URL.String(), initclient.InitApiEndpoint+"/"+orgName+"/systems/"+systemName)

			return nil
		}

		testInitClient := initclient.NewInitClient("")

		testInitClient.R.C.SetTransport(client.RoundTripFunc(mockTransport))

		s, err := testInitClient.GetSystem(orgName, systemName)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestInitClient_GetSystemFromHost(t *testing.T) {
	t.Run("ValidHost", func(tt *testing.T) {
		host := fmt.Sprintf("%s.%s", orgName, systemName)

		testInitClient := initclient.NewInitClient("")

		s, err := testInitClient.GetSystemFromHost(host, nil)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})

	t.Run("InvalidHost", func(tt *testing.T) {
		host := fmt.Sprintf("%s.%s.", orgName, systemName)

		testInitClient := initclient.NewInitClient("")

		s, err := testInitClient.GetSystemFromHost(host, nil)

		assert.Error(tt, err)
		assert.Nil(tt, s)
	})
}

func TestGetHostUrl(t *testing.T) {
	t.Run("InvalidHostScheme", func(tt *testing.T) {
		var org *string = nil

		initclientMock := &cmocks.InitClient{}
		host := fmt.Sprintf("%s.%s.", orgName, systemName)

		initclientMock.On("GetSystemFromHost", host, org).
			Return(nil, errors.New("some error"))

		u, err := initclient.GetHostUrl(initclientMock, host, org)

		assert.Error(tt, err)
		assert.Nil(tt, u)
	})

	t.Run("MissingOrgInHost", func(tt *testing.T) {
		var org *string = nil

		initclientMock := &cmocks.InitClient{}
		host := fmt.Sprintf("%s", systemName)

		initclientMock.On("GetSystemFromHost", host, org).
			Return(&initclient.SystemIPInfo{
				SystemId:   uuid.NewV4().String(),
				SystemName: systemName,
				OrgName:    orgName,
				Ip:         testIp,
				Port:       testPort,
			}, nil)

		u, err := initclient.GetHostUrl(initclientMock, host, nil)

		assert.NoError(tt, err)
		assert.NotNil(tt, u)
	})
}

func TestParseHostString(t *testing.T) {

	t.Run("InvalidHostScheme", func(tt *testing.T) {

	})
}

func TestTestParseHostString(t *testing.T) {
	var org string = "org"

	tests := []struct {
		name  string
		host  string
		org   *string
		isErr bool
	}{
		{
			name:  "3 parts host",
			host:  "org.system.",
			org:   &org,
			isErr: true,
		},

		{
			name:  "host org and org param missmatch",
			host:  "missmatch-org.system",
			org:   &org,
			isErr: true,
		},

		{
			name:  "system name with nil org",
			host:  "system",
			org:   nil,
			isErr: true,
		},

		{
			name:  "system name with non nil org",
			host:  "system",
			org:   &org,
			isErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := initclient.ParseHostString(test.host, test.org)
			assert.Equal(t, test.isErr, err != nil)
		})
	}
}
