package rest

// org group

type GetOrgsRequest struct {
	UserUuid string `example:"{{UserUUID}}" form:"user_uuid" json:"user_uuid" query:"user_uuid" binding:"required" validate:"required"`
}

type GetOrgRequest struct {
	OrgName string `example:"milky-way" path:"org" validate:"required"`
}

type AddOrgRequest struct {
	OrgName     string `example:"milky-way" json:"name" validate:"required"`
	Owner       string `example:"{{UserUUID}}" json:"owner_uuid" validate:"required"`
	Certificate string `example:"test_cert" json:"certificate"`
}

type MemberRequest struct {
	OrgName  string `example:"milky-way" path:"org" validate:"required"`
	UserUuid string `example:"{{UserUUID}}" json:"user_uuid" validate:"required"`
	Role     string `example:"member" json:"role" validate:"required"`
}

type GetMemberRequest struct {
	OrgName  string `example:"milky-way" path:"org" validate:"required"`
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}

type GetMemberRoleRequest struct {
	OrgId    string `example:"{{OrgId}}" path:"org" validate:"required"`
	UserUuid string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
}
type UpdateMemberRequest struct {
	OrgName       string `example:"milky-way" path:"org" validate:"required"`
	UserUuid      string `example:"{{UserUUID}}" path:"user_uuid" validate:"required"`
	IsDeactivated bool   `example:"false" json:"isDeactivated,omitempty"`
	Role          string `example:"member" json:"role,omitempty"`
}

// Users group

type GetUserRequest struct {
	UserId string `example:"{{UserID}}" path:"user_id" validate:"required"`
}

type GetUserByAuthIdRequest struct {
	AuthId string `example:"{{AuthId}}" path:"auth_id" validate:"required"`
}

type AddUserRequest struct {
	Name  string `example:"John" json:"name,omitempty" validate:"required"`
	Email string `example:"john@example.com" json:"email" validate:"required"`
	Phone string `example:"4151231234" json:"phone,omitempty"`
}

// Network group

type GetNetworksRequest struct {
	OrgUuid string `example:"{{OrgUUID}}" form:"org" json:"org" query:"org" binding:"required" validate:"required"`
}

type GetNetworkRequest struct {
	NetworkId string `example:"{{NetworkUUID}}" path:"net_id" validate:"required"`
}

type AddNetworkRequest struct {
	OrgName string `example:"milky-way"  json:"org" validate:"required"`
	NetName string `example:"mesh-network" json:"network_name" validate:"required"`
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
	ParentNode string `json:"parent_node" validate:"required"`
	AmpNodeL   string `json:"anodel" validate:"required"`
	AmpNodeR   string `json:"anoder" validate:"required"`
}

type DetachNodeRequest struct {
	Node string `json:"node" validate:"required"`
}

type UpdateNodeStateRequest struct {
	Node  string `json:"node" path:"node_id" validate:"required"`
	State string `json:"state" path:"state" validate:"required"`
}

type UpdateNodeRequest struct {
	Node string `json:"node" path:"node_id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type GetNodeRequest struct {
	Node string `json:"node" path:"node_id" validate:"required"`
}

type GetAllNodesRequest struct {
}
type GetFreeNodesRequest struct {
}

type AddNodeRequest struct {
	Node  string `json:"node" path:"node_id" validate:"required"`
	State string `json:"state" validate:"required"`
}

type DeleteNodeRequest struct {
	Node string `json:"node" path:"node_id" validate:"required"`
}

type AddNodeToNetworkRequest struct {
	Node    string `json:"node" path:"node_id" validate:"required"`
	Network string `json:"network" path:"network_id" validate:"required"`
}

type ReleaseNodeFromNetwork struct {
	Node string `json:"node" path:"node_id" validate:"required"`
}
