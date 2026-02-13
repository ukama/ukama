/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type GetNodeMetricsInput struct {
	FilterBase
	NodeID string `path:"node" validate:"required"`
}

type GetOrgMetricsInput struct {
	Metric string `path:"metric" validate:"required"`
}

type GetNetworkMetricsInput struct {
	Metric  string `path:"metric" validate:"required"`
	Network string `path:"network" validate:"required"`
}

type FilterBase struct {
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" default:"0"`      // can be omitted, if omitted the Now() is used
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}

type GetSubscriberMetricsInput struct {
	FilterBase
	Metric     string `path:"metric" validate:"required"`
	Network    string `path:"network" validate:"required"`
	Subscriber string `path:"subscriber" validate:"required"`
}

type GetSiteMetricsInput struct {
	FilterBase
	SiteID string `path:"site" validate:"required"`
	Metric string `path:"metric" validate:"required"`
}

type GetMetricsRangeInput struct {
	FilterBase
	Network    string `query:"network"`
	Subscriber string `query:"subscriber" `
	Sim        string `query:"sim"`
	User       string `query:"user"`
	Site       string `query:"site"`
	NodeID     string `query:"node"`
	Operation  string `query:"operation" default:"avg" validate:"oneof=avg sum"`
}
type GetMetricsInput struct {
	Metric string `path:"metric" validate:"required"`
}

type GetSimMetricsInput struct {
	FilterBase
	Metric     string `path:"metric" validate:"required"`
	Network    string `path:"network" validate:"required"`
	Subscriber string `path:"subscriber" validate:"required"`
	Sim        string `path:"sim" validate:"required"`
}

type GetWsMetricInput struct {
	Metric     string `query:"metric" validate:"required"`
	Interval   int    `query:"interval" validate:"required"` //Node/Network/Organization/Site/Sub/User
	Network    string `query:"network"`
	Subscriber string `query:"subscriber"`
	Sim        string `query:"sim"`
	User       string `query:"user"`
	Site       string `query:"site"`
	NodeID     string `query:"node"`
	Operation  string `query:"operation" default:"avg" validate:"oneof=avg sum"`
}

type DummyParameters struct {
}

type GetAlgoStatsForMetricInput struct {
	NodeID     string `path:"node" validate:"required"`
	Metric string `path:"metric" validate:"required"`
}
