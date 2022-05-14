package server

import "github.com/ukama/ukama/testing/services/network/internal/db"

type ReqActionOnNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_to" validate:"eq=vnode_power_on|eq=vnode_power_off,required"`
}

type ReqPowerOnNode struct {
	ReqActionOnNode
}

type ReqPowerOffNode struct {
	ReqActionOnNode
}

type ReqGetNode struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=vnode_info,required"`
}

type RespGetNode struct {
	Node    db.VNode `json:"node"`
	Runtime string   `json:"runtime"`
}

type ReqGetNodeList struct {
	NodeID     string `query:"node" validate:"required"`
	LookingFor string `query:"looking_for" validate:"eq=vnode_list,required"`
}

type RespGetNodeList struct {
	NodeList []db.VNode `json:"nodes"`
}

type ReqDeleteNode struct {
	NodeID    string `query:"node" validate:"required"`
	LookingTo string `query:"looking_for" validate:"eq=vnode_delete,required"`
}
