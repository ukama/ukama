package rest

type DummyParameters struct {
}

type AddOrgRequest struct {
	OrgName     string `path:"org" validate:"required"`
	Ip          string `json:"ip" validate:"required"`
	Certificate string `json:"certificate" validate:"required"`
}

type UpdateOrgRequest struct {
	OrgName     string `path:"org" validate:"required"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
}

type GetOrgRequest struct {
	OrgName string `path:"org" validate:"required"`
}

type GetOrgResponse struct {
	OrgName     string `json:"org"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
}

type AddNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type DeleteNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type GetNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type AddSystemRequest struct {
	OrgName     string `path:"org" validate:"required"`
	SysName     string `path:"system" validate:"required"`
	Ip          string `json:"ip" validate:"required"`
	Certificate string `json:"certificate" validate:"required"`
	Port        int32  `json:"port" validate:"required"`
}

type UpdateSystemRequest struct {
	OrgName     string `path:"org" validate:"required"`
	SysName     string `path:"system" validate:"required"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
	Port        int32  `json:"port"`
}

type AddSystemResponse struct {
	OrgName     string `path:"org"`
	SysName     string `path:"system"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
	Port        int32  `json:"port"`
	Health      int32  `json:"health"`
}

type GetSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}

type GetSystemResponse struct {
	OrgName     string `json:"org"`
	SystemName  string `json:"system"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
	Port        string `json:"port"`
	Health      int32  `json:"health"`
}

type DeleteSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}
