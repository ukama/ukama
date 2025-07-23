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
	UserId   string `example:"{{userId}}" path:"uuid" validate:"required"`
	Category string `example:"{{componentCategory}}" query:"category" validate:"eq=all|eq=access|eq=backhaul|eq=power|eq=switch|eq=spectrum"`
}

type GetAccounts struct {
	UserId string `example:"{{userId}}" path:"uuid" validate:"required"`
}

type GetContracts struct {
	Company  string `example:"{{company}}" path:"company" validate:"required"`
	IsActive bool   `example:"{{true}}" query:"is_active" validate:"required"`
}
