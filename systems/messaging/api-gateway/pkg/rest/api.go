package rest

import (
	"github.com/ukama/ukama/systems/common/rest"
)

type GetNodeIPRequest struct {
	rest.BaseRequest
	NodeId string `path:"node_id" validate:"required"`
}

type SetNodeIPRequest struct {
	rest.BaseRequest
	NodeId   string `path:"node_id" validate:"required"`
	NodeIp   string `json:"node_ip" validate:"required"`
	MeshIp   string `json:"mesh_ip" validate:"required"`
	NodePort int32  `json:"node_port" validate:"required"`
	MeshPort int32  `json:"mesh_port" validate:"required"`
	Org      string `json:"org" validate:"required"`
	Network  string `json:"network" validate:"required"`
}

type DeleteNodeIPRequest struct {
	rest.BaseRequest
	NodeId string `path:"node_id" validate:"required"`
}

type ListNodeIPsRequest struct {
	rest.BaseRequest
}

type NodeOrgMapListRequest struct {
	rest.BaseRequest
}
