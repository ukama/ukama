/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest_test

import "net/http"

const (
	testUuid   = "03cb753f-5e03-4c97-8e47-625115476c72"
	testNodeId = "uk-sa2341-hnode-v0-a1a0"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}
