package rest

import (
	"github.com/ukama/ukama/systems/common/rest"
)

type MemberRequest struct {
	rest.BaseRequest
	UserUuid string `example:"{{UserUUID}}" json:"user_uuid" validate:"required"`
	Role     string `example:"member" json:"role" validate:"required"`
}

type GetMembersRequest struct {
	rest.BaseRequest
}

type GetMemberRequest struct {
	rest.BaseRequest
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}

type RemoveMemberRequest struct {
	rest.BaseRequest
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}

type GetMemberRoleRequest struct {
	rest.BaseRequest
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}
type UpdateMemberRequest struct {
	rest.BaseRequest
	UserUuid      string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
	IsDeactivated bool   `example:"false" json:"isDeactivated,omitempty"`
	Role          string `example:"member" json:"role,omitempty"`
}

// Network group
type GetNetworksRequest struct {
	rest.BaseRequest
	OrgUuid string `example:"{{OrgUUID}}" form:"org" json:"org" query:"org" binding:"required" validate:"required"`
}

type GetNetworkRequest struct {
	rest.BaseRequest
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
}

type AddNetworkRequest struct {
	rest.BaseRequest
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
	rest.BaseRequest
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
	SiteName  string `example:"s1-site" path:"site" validate:"required"`
}
type AddSiteRequest struct {
	rest.BaseRequest
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
	SiteName  string `example:"s1-site" json:"site" validate:"required"`
}

type AttachNodesRequest struct {
	rest.BaseRequest
	ParentNode string `json:"node_id" path:"node_id" validate:"required"`
	AmpNodeL   string `json:"anodel"`
	AmpNodeR   string `json:"anoder"`
}

type DetachNodeRequest struct {
	rest.BaseRequest
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type UpdateNodeStateRequest struct {
	rest.BaseRequest
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
	State  string `json:"state" validate:"required"`
}

type UpdateNodeRequest struct {
	rest.BaseRequest
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type GetNodeRequest struct {
	rest.BaseRequest
	NodeId string `json:"node_id" path:"node_id" validate:"required"`
}

type GetNodesRequest struct {
	rest.BaseRequest
	Free bool `form:"free" json:"free" query:"free" binding:"required"`
}

type GetOrgNodesRequest struct {
	rest.BaseRequest
	OrgId string `example:"{{OrgId}}" path:"org" validate:"required"`
	Free  bool   `form:"free" json:"free" query:"free" binding:"required"`
}

type GetSiteNodesRequest struct {
	rest.BaseRequest
	SiteId string `example:"{{SiteId}}" path:"site_id" validate:"required"`
}

type GetNetworkNodesRequest struct {
	rest.BaseRequest
	NetworkId string `example:"{{NetworkId}}" path:"net_id" validate:"required"`
}

type AddNodeRequest struct {
	rest.BaseRequest
	NodeId string `json:"node_id" validate:"required"`
	Name   string `json:"name"`
	OrgId  string `json:"org_id" validate:"required"`
	State  string `json:"state" validate:"required"`
}

type DeleteNodeRequest struct {
	rest.BaseRequest
	NodeId string `json:"node" path:"node_id" validate:"required"`
}

type AddNodeToSiteRequest struct {
	rest.BaseRequest
	NodeId string `json:"node_id" path:"node_id" validate:"required"`

	// TODO: update RPC handlers for missing site_id (default site for network)
	SiteId    string `json:"site_id"`
	NetworkId string `json:"net_id" validate:"required"`
}

type ReleaseNodeFromSiteRequest struct {
	rest.BaseRequest
	NodeId string `json:"node" path:"node_id" validate:"required"`
}

type AddInvitationRequest struct {
	rest.BaseRequest
	Org   string `json:"org" path:"org" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type GetInvitationByOrgRequest struct {
	rest.BaseRequest
	Org string `json:"org" path:"org" validate:"required"`
}

type GetInvitationRequest struct {
	rest.BaseRequest
	InvitationId string `json:"invitation_id" path:"invitation_id" validate:"required"`
}

type UpdateInvitationRequest struct {
	rest.BaseRequest
	InvitationId string `json:"invitation_id" validate:"required" path:"invitation_id"`
	Status       string `form:"status" json:"status" validate:"required"`
}

type RemoveInvitationRequest struct {
	rest.BaseRequest
	InvitationId string `json:"invitation_id" path:"invitation_id" validate:"required"`
}
