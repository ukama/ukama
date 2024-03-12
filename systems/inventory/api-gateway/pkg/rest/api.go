/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetTestRequest struct{}

type GetComponent struct {
	uuid string `example:"{{ComponentUUID}}" path:"uuid" validate:"required"`
}

type GetComponents struct {
	company       string `example:"{{company}}" path:"uuid" validate:"required"`
	componentType string `example:"{{componentType}}" path:"query" validate:"required"`
}
