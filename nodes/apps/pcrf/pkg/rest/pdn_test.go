/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
)

func TestPDNHTTPRejectsLegacyAndMalformedRequests(t *testing.T) {
	cfg := pkg.NewConfig("pcrf")
	cfg.DebugMode = false
	r := NewRouter(nil, NewRouterConfig(cfg), "test", RuntimeStatus{})
	cases := []struct{ method, path, body string }{
		{"POST", "/v1/session", `{"imsi":[0,0,1,0,1],
            "pdn_address":3232237570,"apn_name":"internet"}`},
		{"DELETE", "/v1/session", `{"imsi":[0,0,1,0,1]}`},
		{"POST", "/v1/session", `{"epc_instance_id":"a","pdn_id":"b",
            "imsi":[0,0,1,0,1],"pdn_address":3232237570}`},
		{"POST", "/v1/session", `{"epc_instance_id":"a","pdn_id":"b",
            "imsi":[0,0,1,0,12],"pdn_address":3232237570,
            "apn_name":"internet"}`},
		{"POST", "/v1/epc", `{}`},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path,
			strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.f.Engine().ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	}
}
