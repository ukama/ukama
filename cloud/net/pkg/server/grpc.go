package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"github.com/ukama/ukamaX/cloud/net/pkg"
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

func (n *NnsServer) Get(c context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	ip, err := n.nns.Get(c, req.GetNodeId())
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Ip: ip,
	}, nil
}

func (n *NnsServer) Set(c context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {

	logrus.Infof("Seting Ip for: %s", req.GetNodeId())

	err := n.nns.Set(c, req.GetNodeId(), req.GetIp())
	if err != nil {
		return nil, fmt.Errorf("failed to add node-ip record to db. Error: %v", err)
	}

	err = n.nodeOrgMapping.Add(c, req.GetNodeId(), req.OrgName, req.Network)
	if err != nil {
		return nil, fmt.Errorf("failed to set org and network for node id %s. Error: %v", req.NodeId, err)
	}

	return &pb.SetResponse{}, nil
}

func (n *NnsServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
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

	return &pb.ListResponse{
		Ips: res,
	}, err
}

func (n *NnsServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := n.nns.Delete(ctx, in.NodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	return &pb.DeleteResponse{}, nil
}
