/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NnsServer struct {
	pb.UnimplementedNnsServer
	config    *pkg.Config
	dnsConfig *pkg.DnsConfig
	nns       *pkg.Nns
}

func NewNnsServer(nnsClient *pkg.Nns, config *pkg.Config, dnsConfig *pkg.DnsConfig) *NnsServer {
	return &NnsServer{
		config:    config,
		dnsConfig: dnsConfig,
		nns:       nnsClient,
	}
}

func (n *NnsServer) GetNode(c context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	log.Infof("Getting node %s", req.GetNodeId())
	if _, err := ukama.ValidateNodeId(req.GetNodeId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	node, err := n.nns.Get(c, req.GetNodeId())
	if err != nil {
		return nil, err
	}
	return &pb.GetNodeResponse{
		NodeId:   node.NodeId,
		NodeIp:   node.NodeIp,
		NodePort: node.NodePort,
	}, nil
}

func (n *NnsServer) GetMesh(c context.Context, req *pb.GetMeshRequest) (*pb.GetMeshResponse, error) {
	log.Infof("Getting mesh")
	mesh, err := n.nns.GetMesh(c)
	if err != nil {
		return nil, err
	}
	return &pb.GetMeshResponse{
		MeshIp:   mesh.MeshIp,
		MeshPort: mesh.MeshPort,
	}, nil
}

func (n *NnsServer) Set(c context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	log.Infof("Setting node %s", req.GetNodeId())
	if _, err := ukama.ValidateNodeId(req.GetNodeId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	obj := pkg.OrgMap{
		NodeId:       req.GetNodeId(),
		NodeIp:       req.GetNodeIp(),
		NodePort:     req.GetNodePort(),
		MeshPort:     req.GetMeshPort(),
		Org:          n.config.OrgName,
		Network:      req.GetNetwork(),
		Site:         req.GetSite(),
		MeshHostName: req.GetMeshHostName(),
		MeshIp:       req.GetMeshIp(),
	}
	log.Infof("Adding node %+v", obj)
	err := n.nns.Add(c, obj)
	if err != nil {
		return nil, err
	}
	return &pb.SetResponse{}, nil
}

func (n *NnsServer) UpdateMesh(c context.Context, req *pb.UpdateMeshRequest) (*pb.UpdateMeshResponse, error) {
	log.Infof("Updating mesh %s:%d", req.GetMeshIp(), req.GetMeshPort())
	err := n.nns.SetMesh(c, req.GetMeshIp(), req.GetMeshPort())
	if err != nil {
		return nil, err
	}
	err = n.nns.UpdateNodeMesh(c, req.GetMeshIp(), req.GetMeshPort())
	if err != nil {
		return nil, err
	}
	return &pb.UpdateMeshResponse{}, nil
}

func (n *NnsServer) UpdateNode(c context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	log.Infof("Updating node %s", req.GetNodeId())
	if _, err := ukama.ValidateNodeId(req.GetNodeId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := n.nns.UpdateNode(c, req.GetNodeId(), req.GetNodeIp(), req.GetNodePort())
	if err != nil {
		return nil, err
	}
	return &pb.UpdateNodeResponse{}, nil
}

func (n *NnsServer) Delete(c context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting node %s", req.GetNodeId())
	if _, err := ukama.ValidateNodeId(req.GetNodeId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := n.nns.Delete(c, req.GetNodeId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{}, nil
}

func (n *NnsServer) List(c context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("Listing all nodes")
	items, err := n.nns.GetAll(c)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListResponse{}
	resp.List = parseOrgMapList(items)
	return resp, nil
}

func parseOrgMap(item pkg.OrgMap) *pb.OrgMap {
	return &pb.OrgMap{
		NodeId:       item.NodeId,
		NodeIp:       item.NodeIp,
		NodePort:     item.NodePort,
		MeshPort:     item.MeshPort,
		Org:          item.Org,
		Network:      item.Network,
		Site:         item.Site,
		MeshIp:       item.MeshIp,
		MeshHostName: item.MeshHostName,
	}
}

func parseOrgMapList(items []pkg.OrgMap) []*pb.OrgMap {
	resp := make([]*pb.OrgMap, 0)
	for _, item := range items {
		resp = append(resp, parseOrgMap(item))
	}
	return resp
}

// func (n *NnsServer) Get(c context.Context, req *pb.GetNodeIPRequest) (*pb.GetNodeIPResponse, error) {
// 	log.Infof("Getting ip for node %s", req.NodeId)
// 	ip, err := n.nns.Get(c, req.GetNodeId())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.GetNodeIPResponse{
// 		Ip: ip,
// 	}, nil
// }

// func (n *NnsServer) Set(c context.Context, req *pb.SetNodeIPRequest) (*pb.SetNodeIPResponse, error) {
// 	log.Infof("Seting Ip for: %s", req.GetNodeId())

// 	i := net.ParseIP(req.GetMeshIp())
// 	if i == nil {
// 		return nil, fmt.Errorf("not valid ip")
// 	}

// 	err := n.nns.Set(c, req.GetNodeId(), fmt.Sprintf("%s:%d", req.GetMeshIp(), req.MeshPort))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to add node-ip record to db. Error: %v", err)
// 	}

// 	err = n.nodeOrgMapping.Add(c, req.GetNodeId(), req.Org, req.Network, req.Site, req.NodeIp, req.MeshHostName, req.NodePort, req.MeshPort)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to set org and network for node id %s. Error: %v", req.NodeId, err)
// 	}

// 	return &pb.SetNodeIPResponse{}, nil
// }

// func (n *NnsServer) List(ctx context.Context, in *pb.ListNodeIPRequest) (*pb.ListNodeIPResponse, error) {
// 	log.Infof("Listing all nodes")
// 	nodes, err := n.nns.List(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get list of nodes. Error: %v", err)
// 	}

// 	ips := make(map[string]interface{})
// 	// makes sure ip list does not have duplicates
// 	for _, v := range nodes {
// 		ips[v] = nil
// 	}
// 	res := make([]string, 0, len(ips))
// 	for k := range ips {
// 		res = append(res, k)
// 	}

// 	return &pb.ListNodeIPResponse{
// 		Ips: res,
// 	}, err
// }

// func (n *NnsServer) Delete(ctx context.Context, in *pb.DeleteNodeIPRequest) (*pb.DeleteNodeIPResponse, error) {
// 	log.Infof("Deleting Ip for: %s", in.GetNodeId())
// 	err := n.nns.Delete(ctx, in.NodeId)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
// 	}

// 	return &pb.DeleteNodeIPResponse{}, nil
// }

// func (n *NnsServer) GetNodeOrgMapList(ctx context.Context, in *pb.NodeOrgMapListRequest) (*pb.NodeOrgMapListResponse, error) {
// 	log.Infof("GetNodeOrgMap List")
// 	resp := &pb.NodeOrgMapListResponse{}
// 	maps, err := n.nodeOrgMapping.List(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
// 	}

// 	for k, v := range maps {
// 		nom := &pb.NodeOrgMap{
// 			NodeId:       k,
// 			NodeIp:       v.NodeIp,
// 			NodePort:     v.NodePort,
// 			MeshPort:     v.MeshPort,
// 			Org:          v.Org,
// 			Network:      v.Network,
// 			Site:         v.Site,
// 			Domainname:   n.config.NodeDomain,
// 			MeshHostName: v.MeshHostName,
// 		}
// 		resp.Map = append(resp.Map, nom)
// 	}
// 	log.Infof("GetNodeOrgMap: %v", resp.Map)
// 	return resp, nil
// }

// func (n *NnsServer) GetNodeIPMapList(ctx context.Context, in *pb.NodeIPMapListRequest) (*pb.NodeIPMapListResponse, error) {
// 	log.Debugf("Node Ip map list request")
// 	resp := &pb.NodeIPMapListResponse{}
// 	m, err := n.nns.List(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to delete record from db. Error: %v", err)
// 	}

// 	for n, i := range m {
// 		e := &pb.NodeIPMap{
// 			NodeId: n,
// 			NodeIp: i,
// 		}
// 		resp.Map = append(resp.Map, e)
// 	}
// 	return resp, nil
// }

// func (n *NnsServer) GetMesh(c context.Context, req *pb.GetMeshIPRequest) (*pb.GetMeshIPResponse, error) {
// 	log.Infof("Getting Mesh IP for node %s", req.NodeId)
// 	orgMap, err := n.nodeOrgMapping.Get(c, req.NodeId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mesh, err := n.nodeOrgMapping.GetMesh(c, orgMap.MeshHostName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pb.GetMeshIPResponse{
// 		Ip:   *mesh,
// 		Port: orgMap.MeshPort,
// 	}, nil
// }

// func (n *NnsServer) SetMesh(ctx context.Context, in *pb.SetMeshRequest) (*pb.SetMeshResponse, error) {
// 	log.Infof("Setting Mesh IP to: %s & port: %d for all records", in.Ip, in.Port)
// 	err := n.nodeOrgMapping.UpdateMesh(ctx, in.Ip, in.Port)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update mesh ip and port for all records. Error: %v", err)
// 	}
// 	return &pb.SetMeshResponse{}, nil
// }
