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
type StoreRunningAppsInfoRequest struct {
	NodeId    string   `validate:"required" example:"{{NodeId}}" path:"node_id" `
	Timestamp string   `json:"timestamp" validate:"required" example:"{{time}}"`
	System    []System `json:"system" validate:"required" example:"{{system}}"`
	Capps     []Capps  `json:"capps" validate:"required" example:"{{capps}}"`
}

type System struct {
	Name  string `json:"name" validate:"required" example:"{{name}}"`
	Value string `json:"value" validate:"required" example:"{{value}}"`
}

type Capps struct {
	Space    string      `json:"space" validate:"required" example:"{{space}}"`
	Name      string      `json:"name" validate:"required" example:"{{name}}"`
	Tag       string      `json:"tag" validate:"required" example:"{{tag}}"`
	Status    string      `json:"status" validate:"required" example:"{{status}}"`
	Resources []Resources `json:"resources" validate:"required" example:"{{resources}}"`
}

type Resources struct {
	Name  string `json:"name" validate:"required" example:"{{name}}"`
	Value string `json:"value" validate:"required" example:"{{value}}"`
}
