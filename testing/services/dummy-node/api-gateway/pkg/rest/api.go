/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type ReqNodeId struct {
	Id string `example:"{{NodeId}}" json:"id" path:"id" validate:"required"`
}

type NodeMetricsById struct {
	Id      string `example:"{{NodeId}}" json:"id" path:"id" validate:"required"`
	Profile string `example:"PROFILE_NORMAL" json:"profile" query:"profile" validate:"eq=PROFILE_NORMAL|eq=PROFILE_MIN|eq=PROFILE_MAX"`
}
