/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPAPNRejection(t *testing.T) {
	cfg := pkg.NewConfig("pcrf")
	cfg.DebugMode = false
	r := NewRouter(nil, NewRouterConfig(cfg), "test", RuntimeStatus{})
	for _, apn := range []interface{}{nil, "", " internet", "a..b", "a\x00b", strings.Repeat("a", 100)} {
		body := map[string]interface{}{"epc_instance_id": "11111111-1111-4111-8111-111111111111",
			"pdn_id": "22222222-2222-4222-8222-222222222222", "imsi": []int{0, 0, 1, 0, 1},
			"pdn_address": 3232237570}
		if apn != nil {
			body["apn_name"] = apn
		}
		data, err := json.Marshal(body)
		require.NoError(t, err)
		req := httptest.NewRequest("POST", "/v1/session", strings.NewReader(string(data)))
		req.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		r.f.Engine().ServeHTTP(response, req)
		require.Equal(t, 400, response.Code)
		require.JSONEq(t, `{"error":"invalid_apn"}`, response.Body.String())
	}
}
