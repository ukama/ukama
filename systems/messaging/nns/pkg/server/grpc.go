package server

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
)

type NnsServer struct {
	pb.UnimplementedNnsServer
	nns            *pkg.Nns
	nodeOrgMapping *pkg.NodeOrgMap
}

func NewNnsServer(nnsClient *pkg.Nns, nodeOrgMapping *pkg.NodeOrgMap) *NnsServer {
	return &NnsServer{
		nns:            nnsClient,
		nodeOrgMapping: nodeOrgMapping,
	}
}

func (n *NnsServer) Get(c context.Context, req *pb.GetNodeIPRequest) (*pb.GetNodeIPResponse, error) {
	logrus.Infof("Getting ip for node %s", req.NodeId)
	ip, err := n.nns.Get(c, req.GetNodeId())
	if err != nil {
		return nil, err
	}
	return &pb.GetNodeIPResponse{
		Ip: ip,
	}, nil
}

func (n *NnsServer) Set(c context.Context, req *pb.SetNodeIPRequest) (*pb.SetNodeIPResponse, error) {
	logrus.Infof("Seting Ip for: %s", req.GetNodeId())

	err := n.nns.Set(c, req.GetNodeId(), req.GetMeshIp())
	if err != nil {
		return nil, fmt.Errorf("failed to add node-ip record to db. Error: %v", err)
	}

	err = n.nodeOrgMapping.Add(c, req.GetNodeId(), req.Org, req.Network, req.NodeIp, req.NodePort, req.MeshPort)
	if err != nil {
		return nil, fmt.Errorf("failed to set org and network for node id %s. Error: %v", req.NodeId, err)
	}

	return &pb.SetNodeIPResponse{}, nil
}

func (n *NnsServer) List(ctx context.Context, in *pb.ListNodeIPRequest) (*pb.ListNodeIPResponse, error) {
	logrus.Infof("Listing all nodes")
	nodes, err := n.nns.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of nodes. Error: %v", err)
	}

	ips := make(map[string]interface{})
	// makes sure ip list does not have duplicates
	for _, v := range nodes {
		ips[v] = nil
	}
	res := make([]string, 0, len(ips))
	for k := range ips {
		res = append(res, k)
	}

	return &pb.ListNodeIPResponse{
		Ips: res,
	}, err
}

func (n *NnsServer) Delete(ctx context.Context, in *pb.DeleteNodeIPRequest) (*pb.DeleteNodeIPResponse, error) {
	logrus.Infof("Deleting Ip for: %s", in.GetNodeId())
	err := n.nns.Delete(ctx, in.NodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	return &pb.DeleteNodeIPResponse{}, nil
}

func (n *NnsServer) GetNodeOrgMapList(ctx context.Context, in *pb.NodeOrgMapListRequest) (*pb.NodeOrgMapListResponse, error) {
	logrus.Infof("GetNodeOrgMap List")
	resp := &pb.NodeOrgMapListResponse{}
	maps, err := n.nodeOrgMapping.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	for k, v := range maps {
		nom := &pb.NodeOrgMap{
			NodeId:   k,
			NodeIp:   v.NodeIp,
			NodePort: v.NodePort,
			MeshPort: v.MeshPort,
			Org:      v.Org,
			Network:  v.Network,
		}
		resp.Map = append(resp.Map, nom)
	}
	logrus.Infof("GetNodeOrgMap: %v", resp.Map)
	return resp, nil
}
