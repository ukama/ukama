package rest

import pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"

type UserRequest struct {
	Org      string `path:"org" validate:"required"`
	SimToken string `json:"simToken"`
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
}

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

type GetUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}

type DeleteUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}

type AddNodeRequest struct {
	OrgName  string `path:"org" validate:"required"`
	NodeId   string `path:"node" validate:"required"`
	NodeName string `json:"name"`
}

func MapNodesList(pbList *pb.GetNodesResponse) *NodesList {
	var nodes []*Node
	for _, node := range pbList.Nodes {
		nodes = append(nodes, &Node{
			NodeId: node.NodeId,
			State:  pb.NodeState_name[int32(node.State)],
			Type:   pb.NodeType_name[int32(node.Type)],
			Name:   node.Name,
		})
	}
	return &NodesList{
		OrgName: pbList.OrgName,
		Nodes:   nodes,
	}
}
