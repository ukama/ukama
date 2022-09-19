package rest

type AddOrgRequest struct {
	OrgName     string `path:"org" validate:"required"`
	Ip          string `json:"ip" validate:"required"`
	Certificate string `json:"certificate" validate:"required"`
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

type AddSystemResponse struct {
	SystemUiid string `json:"uuid"`
}

type GetSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}

type GetSystemResponse struct {
	OrgName     string `json:"org"`
	SystemName  string `json:"system"`
	SystemUiid  string `json:"uuid"`
	Ip          string `json:"ip"`
	Certificate string `json:"certificate"`
	Port        string `json:"port"`
}

type DeleteSystemRequest struct {
	OrgName string `path:"org" validate:"required"`
	SysName string `path:"system" validate:"required"`
}
