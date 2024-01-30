/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

type MemberRequest struct {
	UserUuid string `example:"{{UserUUID}}" json:"user_uuid" validate:"required"`
	Role     string `example:"member" json:"role" validate:"required"`
}

type GetMembersRequest struct {
}

type GetMemberRequest struct {
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}

type RemoveMemberRequest struct {
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}

type GetMemberRoleRequest struct {
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}
type UpdateMemberRequest struct {
	UserUuid      string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
	IsDeactivated bool   `example:"false" json:"isDeactivated,omitempty"`
	Role          string `example:"member" json:"role,omitempty"`
}

// Network group
type GetNetworksRequest struct {
	OrgUuid string `example:"{{OrgUUID}}" form:"org" json:"org" query:"org" binding:"required" validate:"required"`
}

type GetNetworkRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
}

type AddNetworkRequest struct {
	OrgName          string   `example:"milky-way"  json:"org" validate:"required"`
	NetName          string   `example:"mesh-network" json:"network_name" validate:"required"`
	AllowedCountries []string `json:"allowed_countries"`
	AllowedNetworks  []string `json:"allowed_networks"`
	Budget           float64  `json:"budget"`
	Overdraft        float64  `json:"overdraft"`
	TrafficPolicy    uint32   `json:"traffic_policy"`
	PaymentLinks     bool     `example:"true" json:"payment_links"`
}

type GetSiteRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
	SiteName  string `example:"s1-site" path:"site" validate:"required"`
}
type AddSiteRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
	SiteName  string `example:"s1-site" json:"site" validate:"required"`
}

type AttachNodesRequest struct {
	ParentNode string `json:"node_id" path:"node_id" validate:"required"`
	AmpNodeL   string `json:"anodel"`
	AmpNodeR   string `json:"anoder"`
}

type DetachNodeRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type UpdateNodeStateRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
	State  string `json:"state" validate:"required"`
}

type UpdateNodeRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Latitude  float64 `json:"latitude"`   
    Longitude float64 `json:"longitude"` 
}

type GetNodeRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type GetNodesRequest struct {
	Free bool `form:"free" json:"free" query:"free" binding:"required"`
}

type GetOrgNodesRequest struct {
	OrgId string `example:"{{OrgId}}" path:"org" validate:"required"`
	Free  bool   `form:"free" json:"free" query:"free" binding:"required"`
}

type GetSiteNodesRequest struct {
	SiteId string `example:"{{SiteId}}" path:"site_id" validate:"required"`
}

type GetNetworkNodesRequest struct {
	NetworkId string `example:"{{NetworkId}}" path:"net_id" validate:"required"`
}

type AddNodeRequest struct {
	NodeId string `json:"node_id" validate:"required"`
	Name   string `json:"name"`
	OrgId  string `json:"org_id" validate:"required"`
	State  string `json:"state" validate:"required"`
	Latitude  float64 `json:"latitude"`   
    Longitude float64 `json:"longitude"` 
}

type DeleteNodeRequest struct {
	NodeId string `json:"node" path:"node_id" validate:"required"`
}

type AddNodeToSiteRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`

	// TODO: update RPC handlers for missing site_id (default site for network)
	SiteId    string `json:"site_id"`
	NetworkId string `json:"net_id" validate:"required"`
}

type ReleaseNodeFromSiteRequest struct {
	NodeId string `json:"node" path:"node_id" validate:"required"`
}

type AddInvitationRequest struct {
	Org   string `json:"org" path:"org" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type GetInvitationByOrgRequest struct {
	Org string `json:"org" path:"org" validate:"required"`
}

type GetInvitationRequest struct {
	InvitationId string `json:"invitation_id" path:"invitation_id" validate:"required"`
}

type UpdateInvitationRequest struct {
	InvitationId string `json:"invitation_id" validate:"required" path:"invitation_id"`
	Status       string `form:"status" json:"status" validate:"required"`
}

type RemoveInvitationRequest struct {
	InvitationId string `json:"invitation_id" path:"invitation_id" validate:"required"`
}
