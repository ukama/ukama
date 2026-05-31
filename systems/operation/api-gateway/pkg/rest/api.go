/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

type StartOperationRequest struct {
	Type           string `json:"type" validate:"required" example:"RestartNode"`
	System         string `json:"system" validate:"required" example:"node"`
	ResourceKey    string `json:"resource_key" validate:"required" example:"node:uk-sa2450-tnode-v0-4e86"`
	RequestedBy    string `json:"requested_by"`
	IdempotencyKey string `json:"idempotency_key"`
	LeaseSeconds   uint32 `json:"lease_seconds"`
}

type GetOperationRequest struct {
	Id string `json:"id" path:"id" validate:"required,uuid"`
}

type GetByResourceRequest struct {
	ResourceKey string `json:"resource_key" query:"resource_key" validate:"required"`
}

type MarkRunningRequest struct {
	Id           string `json:"id" path:"id" validate:"required,uuid"`
	FencingToken uint64 `json:"fencing_token" validate:"required"`
}

type ForceUnlockRequest struct {
	Id     string `json:"id" path:"id" validate:"required,uuid"`
	Actor  string `json:"actor" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}
