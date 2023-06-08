package rest

type GetNodeIPRequest struct {
	NodeId string `path:"node" validate:"required"`
}

type SetNodeIPRequest struct {
	NodeId   string `path:"node" validate:"required"`
	NodeIp   string `json:"node_ip" validate:"required"`
	MeshIp   string `json:"mesh_ip" validate:"required"`
	NodePort int32  `json:"node_port" validate:"required"`
	MeshPort int32  `json:"mesh_port" validate:"required"`
	Org      string `json:"org" validate:"required"`
	Network  string `json:"network" validate:"required"`
}

type DeleteNodeIPRequest struct {
	NodeId string `path:"node" validate:"required"`
}

type ListNodeIPsRequest struct {
}

type NodeOrgMapListRequest struct {
}
