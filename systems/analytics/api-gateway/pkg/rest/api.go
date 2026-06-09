/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

/* Business */

type GetBusinessHomeRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	Period    string `example:"week" json:"period" query:"period"`
	From      string `example:"2024-01-01T00:00:00Z" json:"from" query:"from"`
	To        string `example:"2024-01-31T23:59:59Z" json:"to" query:"to"`
	Timezone  string `example:"UTC" json:"timezone" query:"timezone"`
}

type GetSalesOverviewRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	Period    string `example:"week" json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type GetPackagePerformanceRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	Period    string `example:"month" json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetBillingSummaryRequest struct {
	Period   string `example:"month" json:"period" query:"period"`
	From     string `json:"from" query:"from"`
	To       string `json:"to" query:"to"`
	Timezone string `json:"timezone" query:"timezone"`
	Page     uint32 `json:"page" query:"page"`
	PageSize uint32 `json:"page_size" query:"page_size"`
}

type GetBusinessSitesRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetBusinessSiteRequest struct {
	SiteId   string `example:"{{SiteUUID}}" path:"site_id" validate:"required"`
	Period   string `json:"period" query:"period"`
	From     string `json:"from" query:"from"`
	To       string `json:"to" query:"to"`
	Timezone string `json:"timezone" query:"timezone"`
}

type GetInventoryReadinessRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
}

/* Customer */

type GetCustomerOverviewRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type ListCustomersRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	SiteId    string `example:"{{SiteUUID}}" json:"site_id" query:"site_id"`
	Status    string `example:"active" json:"status" query:"status"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type SearchCustomersRequest struct {
	Query     string `example:"john" json:"q" query:"q" validate:"required"`
	NetworkId string `json:"network_id" query:"network_id"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetCustomerRequest struct {
	CustomerId string `example:"{{CustomerUUID}}" path:"customer_id" validate:"required"`
	Period     string `json:"period" query:"period"`
	From       string `json:"from" query:"from"`
	To         string `json:"to" query:"to"`
	Timezone   string `json:"timezone" query:"timezone"`
}

type GetCustomerSupportRequest struct {
	CustomerId string `example:"{{CustomerUUID}}" path:"customer_id" validate:"required"`
}

type GetSimsRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	Status    string `json:"status" query:"status"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetSimPoolRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
}

/* Network */

type GetNetworkOverviewRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type GetTopologyRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
}

type GetNetworkSitesRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	Status    string `example:"online" json:"status" query:"status"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetNetworkSiteRequest struct {
	SiteId   string `example:"{{SiteUUID}}" path:"site_id" validate:"required"`
	Period   string `json:"period" query:"period"`
	From     string `json:"from" query:"from"`
	To       string `json:"to" query:"to"`
	Timezone string `json:"timezone" query:"timezone"`
}

type GetNetworkNodesRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	Status    string `json:"status" query:"status"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetNetworkNodeRequest struct {
	NodeId   string `example:"{{NodeId}}" path:"node_id" validate:"required"`
	Period   string `json:"period" query:"period"`
	From     string `json:"from" query:"from"`
	To       string `json:"to" query:"to"`
	Timezone string `json:"timezone" query:"timezone"`
}

type GetNodePoolRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
}

type GetRadioRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	NodeId    string `json:"node_id" query:"node_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type GetBackhaulRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type GetPowerRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type GetAlarmsRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	Severity  string `example:"critical" json:"severity" query:"severity"`
	State     string `example:"open" json:"state" query:"state"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type GetMetricsRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	NodeId    string `json:"node_id" query:"node_id"`
	Metric    string `example:"uptime" json:"metric" query:"metric"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
}

type GetEventsRequest struct {
	NetworkId string `json:"network_id" query:"network_id"`
	SiteId    string `json:"site_id" query:"site_id"`
	NodeId    string `json:"node_id" query:"node_id"`
	Period    string `json:"period" query:"period"`
	From      string `json:"from" query:"from"`
	To        string `json:"to" query:"to"`
	Timezone  string `json:"timezone" query:"timezone"`
	Page      uint32 `json:"page" query:"page"`
	PageSize  uint32 `json:"page_size" query:"page_size"`
}

type SupportSearchRequest struct {
	Query     string `example:"site-a" json:"q" query:"q" validate:"required"`
	NetworkId string `json:"network_id" query:"network_id"`
}

/* Collector */

type PostRefreshRequest struct {
	Source string `example:"all" json:"source" validate:"required"`
}

type GetRefreshStateRequest struct {
}

type PostRebuildRollupsRequest struct {
	Family string `example:"all" json:"family" validate:"required"`
	From   string `example:"2024-01-01T00:00:00Z" json:"from"`
	To     string `example:"2024-01-31T23:59:59Z" json:"to"`
}

type PostSeedDemoRequest struct {
	Sites     uint32 `json:"sites"`
	Nodes     uint32 `json:"nodes"`
	Customers uint32 `json:"customers"`
	Days      uint32 `json:"days"`
}
