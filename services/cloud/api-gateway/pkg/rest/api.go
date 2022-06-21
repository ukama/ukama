package rest

import (
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	pbnode "github.com/ukama/ukama/services/cloud/node/pb/gen"
)

// Users group

type GetUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}

type UserRequest struct {
	Org      string `path:"org" validate:"required"`
	SimToken string `json:"simToken"`
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	OrgName       string `path:"org" validate:"required"`
	UserId        string `path:"user" validate:"required"`
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	IsDeactivated bool   `json:"isDeactivated,omitempty"`
}

type SetSimStatusRequest struct {
	OrgName string       `path:"org" validate:"required"`
	UserId  string       `path:"user" validate:"required"`
	Iccid   string       `path:"iccid" validate:"required"`
	Carrier *SimServices `json:"carrier,omitempty"`
	Ukama   *SimServices `json:"ukama,omitempty"`
}

type GetSimQrRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
	Iccid   string `path:"iccid" validate:"required"`
}

type SimServices struct {
	Voice *bool `json:"voice,omitempty"`
	Sms   *bool `json:"sms,omitempty"`
	Data  *bool `json:"data,omitempty"`
}

type DeleteUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}

// Nodes group

type NodesList struct {
	OrgName string  `json:"orgName"`
	Nodes   []*Node `json:"nodes"`
}

type Node struct {
	NodeId string `json:"nodeId,omitempty"`
	State  string `json:"state,omitempty"`
	Type   string `json:"type,omitempty"`
	Name   string `json:"name"`
}

type NodeExtended struct {
	Node
	Attached []*Node `json:"attached,omitempty"`
}

type GetNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

type AddNodeRequest struct {
	OrgName  string `path:"org" validate:"required"`
	NodeId   string `path:"node" validate:"required"`
	NodeName string `json:"name" validate:"max=255"`
}

type DeleteNodeRequest struct {
	OrgName string `path:"org" validate:"required"`
	NodeId  string `path:"node" validate:"required"`
}

func MapNodesList(pbList *pb.GetNodesResponse) *NodesList {
	var nodes []*Node
	for _, node := range pbList.Nodes {
		nodes = append(nodes, mapPbNode(node))
	}
	return &NodesList{
		OrgName: pbList.OrgName,
		Nodes:   nodes,
	}
}

func mapPbNode(node *pb.Node) *Node {
	return &Node{
		NodeId: node.NodeId,
		State:  pb.NodeState_name[int32(node.State)],
		Type:   pb.NodeType_name[int32(node.Type)],
		Name:   node.Name,
	}
}

func mapExtendeNode(node *pbnode.Node) *NodeExtended {
	nx := &NodeExtended{
		Node: Node{
			NodeId: node.NodeId,
			State:  pb.NodeState_name[int32(node.State)],
			Type:   pb.NodeType_name[int32(node.Type)],
			Name:   node.Name,
		},
	}
	if len(node.Attached) > 0 {
		nx.Attached = make([]*Node, len(node.Attached))
		for i, n := range node.Attached {
			nx.Attached[i] = mapNodePbNode(n)
		}
	}
	return nx
}

func mapNodePbNode(node *pbnode.Node) *Node {
	return &Node{
		NodeId: node.NodeId,
		State:  pb.NodeState_name[int32(node.State)],
		Type:   pb.NodeType_name[int32(node.Type)],
		Name:   node.Name,
	}
}
