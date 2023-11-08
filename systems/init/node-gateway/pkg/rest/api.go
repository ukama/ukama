/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetNodeRequest struct {
	NodeId string `path:"node" validate:"required"`
}

type GetNodeResponse struct {
	NodeId      string `path:"node" validate:"required"`
	OrgName     string `path:"org" validate:"required"`
	Certificate string `json:"certificate"`
	Ip          string `json:"ip"`
}
