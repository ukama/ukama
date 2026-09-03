/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/common/rest"
)

func newTestRemoteControllerClient(t *testing.T, srv *httptest.Server) *remoteControllerClient {
	t.Helper()

	u, err := url.Parse(srv.URL)
	assert.NoError(t, err)

	return &remoteControllerClient{
		u: u,
		R: rest.NewRestyClient(u, false),
	}
}

func TestGetSubscriberProfile_NotFound_ReturnsErrRemoteSubscriberNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := newTestRemoteControllerClient(t, srv)

	_, err := c.GetSubscriberProfile("999991000000099")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrRemoteSubscriberNotFound))
}

func TestGetSubscriberProfile_OtherErrorStatus_ReturnsPlainError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newTestRemoteControllerClient(t, srv)

	_, err := c.GetSubscriberProfile("999991000000099")
	assert.Error(t, err)
	assert.False(t, errors.Is(err, ErrRemoteSubscriberNotFound))
}

func TestGetSubscriberProfile_Success_ReturnsProfile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"imsi":"999991000000099","reroute":"http://localhost:8085"}`))
	}))
	defer srv.Close()

	c := newTestRemoteControllerClient(t, srv)

	spr, err := c.GetSubscriberProfile("999991000000099")
	assert.NoError(t, err)
	assert.Equal(t, "999991000000099", spr.Imsi)
}
