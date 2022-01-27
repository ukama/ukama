package rest

import pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"

type UserRequest struct {
	Org       string `path:"org" validate:"required"`
	Imsi      string `json:"imsi" validate:"required"`
	FirstName string `json:"firstName,omitempty" validate:"required"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

type GetNodeMetricsInput struct {
	Org    string `path:"org" validate:"required"`
	NodeID string `path:"node" validate:"required"`
	Metric string `path:"metric" validate:"required"`
	From   int64  `query:"from" validate:"required"`
	To     int64  `query:"to" validate:"required"`
	Step   uint   `query:"step" default:"3600"` // default 1 hour
}

type NodesList struct {
	OrgName string  `json:"orgName"`
	Nodes   []*Node `json:"nodes"`
}

type Node struct {
	NodeId string `json:"nodeId,omitempty"`
	State  string `json:"state,omitempty"`
}

func MapNodesList(pbList *pb.NodesList) *NodesList {
	var nodes []*Node
	for _, node := range pbList.Nodes {
		nodes = append(nodes, &Node{
			NodeId: node.NodeId,
			State:  pb.NodeState_name[int32(node.State)],
		})
	}
	return &NodesList{
		OrgName: pbList.OrgName,
		Nodes:   nodes,
	}
}
