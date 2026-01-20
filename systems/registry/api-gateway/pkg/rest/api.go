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
	MemberId string `example:"{{MemberId}}" path:"member_id" validate:"required"`
}

type GetMemberByUserRequest struct {
	UserId string `example:"{{UserId}}" path:"user_id" validate:"required"`
}

type RemoveMemberRequest struct {
	MemberId string `example:"{{MemberId}}" path:"member_id" validate:"required"`
}

type GetMemberRoleRequest struct {
	MemberId string `example:"{{MemberId}}" path:"member_id" validate:"required"`
}
type UpdateMemberRequest struct {
	MemberId      string `example:"{{MemberId}}" path:"member_id" validate:"required"`
	IsDeactivated bool   `example:"false" json:"isDeactivated,omitempty"`
	Role          string `example:"member" json:"role,omitempty"`
}

type GetNetworkRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
}

type AddNetworkRequest struct {
	NetName          string   `example:"mesh-network" json:"network_name" validate:"required"`
	AllowedCountries []string `json:"allowed_countries"`
	AllowedNetworks  []string `json:"allowed_networks"`
	Budget           float64  `json:"budget"`
	Overdraft        float64  `json:"overdraft"`
	TrafficPolicy    uint32   `json:"traffic_policy"`
	PaymentLinks     bool     `example:"true" json:"payment_links"`
}

type GetSiteRequest struct {
	SiteId string `example:"{{SiteUUID}}" path:"site_id" validate:"required"`
}

type GetSitesListRequest struct {
	NetworkId     string `example:"{{NetworkUUID}}" json:"network_id" query:"network_id" binding:"required"`
	IsDeactivated bool   `json:"is_deactivated" query:"is_deactivated" binding:"required" default:"false"`
}
type UpdateSiteRequest struct {
	SiteId string `example:"{{SiteUUID}}" path:"site_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type AddSiteRequest struct {
	NetworkId     string  `example:"{{NetworkUUID}}" json:"network_id" validate:"required"`
	Name          string  `example:"s1-site" json:"site" validate:"required"`
	Location      string  `example:"location" json:"location" validate:"required"`
	BackhaulId    string  `example:"{{BackhaulUUID}}" json:"backhaul_id" validate:"required"`
	PowerId       string  `example:"{{PowerUUID}}" json:"power_id" validate:"required"`
	AccessId      string  `example:"{{AccessUUID}}" json:"access_id" validate:"required"`
	SwitchId      string  `example:"{{SwitchUUID}}" json:"switch_id" validate:"required"`
	SpectrumId    string  `example:"{{SpectrumUUID}}" json:"spectrum_id" validate:"required"`
	IsDeactivated bool    `json:"is_deactivated"`
	Latitude      string  `json:"latitude" validate:"required"`
	Longitude     string  `json:"longitude" validate:"required"`
	InstallDate   string  `json:"install_date"`
}

type AttachNodesRequest struct {
	ParentNode string `json:"node_id" path:"node_id" validate:"required"`
	AmpNodeL   string `json:"anodel"`
	AmpNodeR   string `json:"anoder"`
}

type DetachNodeRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}
type ListNodesRequest struct {
	NetworkId    string `json:"network_id" query:"network_id"`
	SiteId       string `json:"site_id" query:"site_id"`
	Connectivity string `json:"connectivity,omitempty" validate:"omitempty,eq=unknown|eq=online|eq=high|eq=offline" query:"connectivity" enums:"unknown,online,high,offline"`
	State        string `json:"state,omitempty" validate:"omitempty,eq=unknown|eq=configured|eq=operational|eq=faulty" query:"state" enums:"unknown,configured,operational,faulty"`
	NodeId       string `json:"node_id" query:"node_id"`
	Type         string `json:"type" query:"type"`
}

type UpdateNodeStateRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
	State  string `json:"state" validate:"required"`
}

type UpdateNodeRequest struct {
	NodeId    string  `json:"node_id" path:"node_id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
}

type GetNodeRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type GetNodesByStateRequest struct {
	Connectivity string `form:"connectivity" json:"connectivity,omitempty" validate:"eq=unknown|eq=online|eq=high|eq=offline" query:"connectivity" binding:"required"`
	State        string `form:"state" json:"state,omitempty" validate:"eq=unknown|eq=configured|eq=operational|eq=faulty" query:"state" binding:"required"`
}

type GetOrgNodesRequest struct {
	Free bool `form:"free" json:"free" query:"free" binding:"required"`
}

type GetSiteNodesRequest struct {
	SiteId string `example:"{{SiteId}}" path:"site_id" validate:"required"`
}

type GetNetworkNodesRequest struct {
	NetworkId string `example:"{{NetworkId}}" path:"net_id" validate:"required"`
}

type AddNodeRequest struct {
	NodeId    string  `json:"node_id" validate:"required"`
	Name      string  `json:"name"`
	State     string  `json:"state" validate:"required"`
}

type DeleteNodeRequest struct {
	NodeId string `json:"node" path:"node_id" validate:"required"`
}

type AddNodeToSiteRequest struct {
	NodeId string `json:"node_id" path:"node_id" validate:"required"`

	// TODO: update RPC handlers for missing site_id (default site for network)
	SiteId    string `json:"site_id"`
	NetworkId string `json:"network_id" validate:"required"`
}

type ReleaseNodeFromSiteRequest struct {
	NodeId string `json:"node" path:"node_id" validate:"required"`
}

type AddInvitationRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type GetInvitationRequest struct {
	InvitationId string `json:"invitation_id" path:"invitation_id" validate:"required"`
}

type UpdateInvitationRequest struct {
	InvitationId string `json:"invitation_id" validate:"required" path:"invitation_id"`
	Email        string `form:"email" json:"email" validate:"required"`
	Status       string `form:"status" json:"status" validate:"required"`
}

type RemoveInvitationRequest struct {
	InvitationId string `json:"invitation_id" path:"invitation_id" validate:"required"`
}

type GetInvitationsByEmailReq struct {
	Email string `json:"email" path:"email" validate:"required"`
}
