package rest

// org group

type GetOrgsRequest struct {
	UserUUID string `form:"user_uuid" json:"user_uuid" query:"user_uuid" binding:"required" validate:"required"`
}

type GetOrgRequest struct {
	OrgName string `path:"org" validate:"required"`
}

type AddOrgRequest struct {
	OrgName     string `json:"name" validate:"required"`
	Owner       string `json:"owner_uuid" validate:"required"`
	Certificate string `json:"certificate"`
}

type MemberRequest struct {
	OrgName  string `path:"org" validate:"required"`
	UserUUID string `json:"user_uuid" validate:"required"`
}

type GetMemberRequest struct {
	OrgName  string `path:"org" validate:"required"`
	UserUUID string `path:"user_uuid" validate:"required"`
}

type UpdateMemberRequest struct {
	OrgName       string `path:"org" validate:"required"`
	UserUUID      string `path:"user_uuid" validate:"required"`
	IsDeactivated bool   `json:"isDeactivated,omitempty"`
}

// Users group

type GetUserRequest struct {
	UserUUID string `path:"user_uuid" validate:"required"`
}

type AddUserRequest struct {
	Name  string `json:"name,omitempty" validate:"required"`
	Email string `json:"email" validate:"required"`
	Phone string `json:"phone,omitempty" validate:"required"`
}

// Network group

type GetNetworksRequest struct {
	OrgName string `form:"org" json:"org" query:"org" binding:"required" validate:"required"`
}

type GetNetworkRequest struct {
	NetworkID string `path:"net_id" validate:"required"`
}

type AddNetworkRequest struct {
	OrgName string `json:"org" validate:"required"`
	NetName string `json:"network_name" validate:"required"`
}

type GetSiteRequest struct {
	NetworkID string `path:"net_id" validate:"required"`
	SiteName  string `path:"site" validate:"required"`
}
type AddSiteRequest struct {
	NetworkID string `path:"net_id" validate:"required"`
	SiteName  string `json:"site" validate:"required"`
}
