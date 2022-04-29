package server

type ReqGetNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=fact_node_info,required"`
}

type RespGetNode struct {
	NodeID string `json:"node"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type ReqBuildNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=fact_update,required"`
	Type      string `query:"type" validate:"eq=HNODE|eq=TNODE|ANODE,required"`
	Count     int    `query:"count" default:1`
}

type ReqDeleteNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=fact_delete,required"`
}

type ReqGetNodeList struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=fact_node_list,required"`
}

type RespGetNodeList struct {
	NodeList []RespGetNode `json: "nodes"`
}