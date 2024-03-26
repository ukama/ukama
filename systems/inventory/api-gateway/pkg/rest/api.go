/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetTestRequest struct{}

type GetRequest struct {
	Uuid string `example:"{{ComponentUUID}}" path:"uuid" validate:"required"`
}

type GetComponents struct {
	Company  string `example:"{{company}}" path:"company" validate:"required"`
	Category string `example:"{{componentType}}" query:"category" validate:"required"`
}

type GetAccounts struct {
	Company string `example:"{{company}}" path:"company" validate:"required"`
}

type GetContracts struct {
	Company  string `example:"{{company}}" path:"company" validate:"required"`
	IsActive bool   `example:"{{true}}" query:"is_active" validate:"required"`
}
