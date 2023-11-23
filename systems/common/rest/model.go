/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}