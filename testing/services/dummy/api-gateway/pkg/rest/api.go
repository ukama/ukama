/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type UpdateReq struct {
	Iccid   string `form:"iccid" json:"iccid" validate:"required"`
	NodeId  string `form:"node_id" json:"node_id,omitempty"`
	Profile string `form:"profile" json:"profile,omitempty" validate:"eq=normal|eq=min|eq=max" enum:"normal,min,max"`
	Status  string `form:"status" json:"status,omitempty" validate:"eq=ACTIVE|eq=INACTIVE" enum:"ACTIVE,INACTIVE"`
}
