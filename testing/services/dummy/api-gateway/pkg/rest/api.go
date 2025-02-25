/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetRunningAppsRequest struct {
	NodeId string `example:"{{NodeId}}" validate:"required" path:"node_id" `
}

type UpdateReq struct {
	SubscriberID string `form:"subscriber_id" json:"subscriber_id"`
	Profile      string `form:"profile" json:"profile,omitempty" validate:"eq=normal|eq=min|eq=max"`
	Scenario     string `form:"type" json:"type,omitempty"`
}
