/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import "time"

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
	UserId string `json:"user_id" validate:"required,uuid"`
	Reason string `json:"reason" validate:"required"`
}

/*
 * REST response DTOs.
 *
 * Keep external REST JSON consistent with the rest of Ukama APIs:
 * snake_case at the HTTP boundary.
 *
 * Do not return protobuf messages directly from REST handlers because
 * protobuf JSON tags expose lowerCamel fields such as:
 *
 *   fencingToken
 *   resourceKey
 *   leaseExpiresAt
 *
 * REST clients should instead see:
 *
 *   fencing_token
 *   resource_key
 *   lease_expires_at
 */

type Operation struct {
	Id             string     `json:"id"`
	Type           string     `json:"type"`
	System         string     `json:"system"`
	Status         string     `json:"status"`
	FencingToken   uint64     `json:"fencing_token"`
	RequestedBy    string     `json:"requested_by"`
	IdempotencyKey string     `json:"idempotency_key"`
	ResourceKey    string     `json:"resource_key"`
	LeaseExpiresAt *time.Time `json:"lease_expires_at"`
	Error          string     `json:"error"`
	StartedAt      *time.Time `json:"started_at"`
	TerminalAt     *time.Time `json:"terminal_at"`
	CreatedAt      *time.Time `json:"created_at"`
}

type StartOperationResponse struct {
	Operation            *Operation `json:"operation"`
	ConflictingOperation *Operation `json:"conflicting_operation,omitempty"`
}

type GetOperationResponse struct {
	Operation *Operation `json:"operation"`
}

type GetByResourceResponse struct {
	Locked    bool       `json:"locked"`
	Operation *Operation `json:"operation"`
}

type MarkRunningResponse struct {
	Operation *Operation `json:"operation"`
}

type ForceUnlockResponse struct {
	Operation *Operation `json:"operation"`
}
