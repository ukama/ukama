/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type UpdateReq struct {
	SiteId  string `form:"siteId" json:"siteId,omitempty"`
	Profile string `form:"profile" json:"profile,omitempty" validate:"required"`
	Scenario string `form:"scenario" json:"scenario,omitempty" validate:"required"`
	PortUpdates []PortUpdate `form:"portUpdates" json:"portUpdates,omitempty" validate:"required"`
}

type PortUpdate struct {
	PortNumber int32 `form:"portNumber" json:"portNumber" validate:"required"`
	Status bool `form:"status" json:"status" validate:"required"`
}