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

type ListComponentsReq struct {
	Id         string `form:"id" json:"id" query:"id" binding:"required"`
	UserId     string `form:"user_id" json:"user_id" query:"user_id" binding:"required"`
	PartNumber string `form:"part_number" json:"part_number" query:"part_number" binding:"required"`
	Category   string `form:"category" default:"all" json:"category" query:"category" binding:"required" validate:"eq=all|eq=access|eq=backhaul|eq=power|eq=switch|eq=spectrum"`
}

type VerifyRequest struct {
	PartNumber string `example:"{{ComponentPartNumber}}" path:"part_number" validate:"required"`
}

type StartSchedulerRequest struct{}

type StopSchedulerRequest struct{}
