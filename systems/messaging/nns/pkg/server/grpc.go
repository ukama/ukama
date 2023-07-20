package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
)

type NnsServer struct {
	pb.UnimplementedNnsServer
	config         *pkg.DnsConfig
	nns            *pkg.Nns
	nodeOrgMapping *pkg.NodeOrgMap
}

func NewNnsServer(nnsClient *pkg.Nns, nodeOrgMapping *pkg.NodeOrgMap, config *pkg.DnsConfig) *NnsServer {
	return &NnsServer{
		nns:            nnsClient,
		nodeOrgMapping: nodeOrgMapping,
		config:         config,
	}
}

func (n *NnsServer) Get(c context.Context, req *pb.GetNodeIPRequest) (*pb.GetNodeIPResponse, error) {
	log.Infof("Getting ip for node %s", req.NodeId)
	ip, err := n.nns.Get(c, req.GetNodeId())
	if err != nil {
		return nil, err
	}
	return &pb.GetNodeIPResponse{
		Ip: ip,
	}, nil
}

func (n *NnsServer) Set(c context.Context, req *pb.SetNodeIPRequest) (*pb.SetNodeIPResponse, error) {
	log.Infof("Seting Ip for: %s", req.GetNodeId())

	err := n.nns.Set(c, req.GetNodeId(), req.GetMeshIp())
	if err != nil {
		return nil, fmt.Errorf("failed to add node-ip record to db. Error: %v", err)
	}

	err = n.nodeOrgMapping.Add(c, req.GetNodeId(), req.Org, req.Network, req.Site, req.NodeIp, req.NodePort, req.MeshPort)
	if err != nil {
		return nil, fmt.Errorf("failed to set org and network for node id %s. Error: %v", req.NodeId, err)
	}

	return &pb.SetNodeIPResponse{}, nil
}

func (n *NnsServer) List(ctx context.Context, in *pb.ListNodeIPRequest) (*pb.ListNodeIPResponse, error) {
	log.Infof("Listing all nodes")
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
	log.Infof("Deleting Ip for: %s", in.GetNodeId())
	err := n.nns.Delete(ctx, in.NodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	return &pb.DeleteNodeIPResponse{}, nil
}

func (n *NnsServer) GetNodeOrgMapList(ctx context.Context, in *pb.NodeOrgMapListRequest) (*pb.NodeOrgMapListResponse, error) {
	log.Infof("GetNodeOrgMap List")
	resp := &pb.NodeOrgMapListResponse{}
	maps, err := n.nodeOrgMapping.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	for k, v := range maps {
		nom := &pb.NodeOrgMap{
			NodeId:     k,
			NodeIp:     v.NodeIp,
			NodePort:   v.NodePort,
			MeshPort:   v.MeshPort,
			Org:        v.Org,
			Network:    v.Network,
			Site:       v.Site,
			Domainname: n.config.NodeDomain,
		}
		resp.Map = append(resp.Map, nom)
	}
	log.Infof("GetNodeOrgMap: %v", resp.Map)
	return resp, nil
}

func (n *NnsServer) GetNodeIPMapList(ctx context.Context, in *pb.NodeIPMapListRequest) (*pb.NodeIPMapListResponse, error) {
	log.Debugf("Node Ip map list request")
	resp := &pb.NodeIPMapListResponse{}
	m, err := n.nns.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
	}

	for n, i := range m {
		e := &pb.NodeIPMap{
			NodeId: n,
			NodeIp: i,
		}
		resp.Map = append(resp.Map, e)
	}
	return resp, nil
}
