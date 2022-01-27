package server

import (
	"context"
	"fmt"
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
)

type NnsServer struct {
	pb.UnimplementedNnsServer
	nnsClient *Nns
}

func NewNnsServer(nnsClient *Nns) *NnsServer {
	return &NnsServer{
		nnsClient: nnsClient,
	}
}

func (n *NnsServer) Get(c context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	ip, err := n.nnsClient.Get(c, req.GetNodeId())
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		Ip: ip,
	}, nil
}

func (n *NnsServer) Set(c context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	err := n.nnsClient.Set(c, req.GetNodeId(), req.GetIp())
	if err != nil {
		return nil, fmt.Errorf("failed to add record to db. Error: %v", err)
	}

	return &pb.SetResponse{}, nil
}

func (n *NnsServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	nodes, err := n.nnsClient.List(ctx)
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
	err := n.nnsClient.Delete(ctx, in.NodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	return &pb.DeleteResponse{}, nil
}
