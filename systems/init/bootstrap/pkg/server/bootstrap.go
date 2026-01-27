/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package server

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	factory "github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/init/bootstrap/client"
	"github.com/ukama/ukama/systems/init/bootstrap/utils"

	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/db"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)
 
const MessagingSystem = "messaging"
 
type BootstrapServer struct {
	pb.UnimplementedBootstrapServiceServer
	bootstrapRoutingKey msgbus.RoutingKeyBuilder
	nodeRepo            db.NodeRepo
	msgbus              mb.MsgBusServiceClient
	debug               bool
	lookupClient        client.LookupClientProvider
	factoryClient       factory.NodeFactoryClient
	dnsMap              map[string]string
	clientSet 			*kubernetes.Clientset
	config 				*pkg.Config
}
 
func NewBootstrapServer(nodeRepo db.NodeRepo, msgBus mb.MsgBusServiceClient, debug bool, lookupClient client.LookupClientProvider, factoryClient factory.NodeFactoryClient, dnsMap map[string]string, config *pkg.Config) *BootstrapServer {
	c, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	cs, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err.Error())
	}
	return &BootstrapServer{
		clientSet:           cs,
		nodeRepo:            nodeRepo,
		bootstrapRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(config.OrgName).SetService(pkg.ServiceName),
		msgbus:              msgBus,
		debug:               debug,
		lookupClient:        lookupClient,
		factoryClient:       factoryClient,
		dnsMap:              dnsMap,
		config:              config,
	}
}
 
func (s *BootstrapServer) GetNodeMeshInfo(ctx context.Context, req *pb.GetNodeMeshInfoRequest) (*pb.GetNodeMeshInfoResponse, error) {
	node, err := s.nodeRepo.GetNode(req.Id)
	if err != nil {
		log.Errorf("failed to get node from database: %v", err)
		return nil, err
	}

	return &pb.GetNodeMeshInfoResponse{
		NodeId: node.NodeId,
		MeshPodName: node.MeshPodName,
		MeshPodIp: node.MeshPodIp,
		MeshPodPort: node.MeshPodPort,
	}, nil
}

func (s *BootstrapServer) GetNodeCredentials(ctx context.Context, req *pb.GetNodeCredentialsRequest) (*pb.GetNodeCredentialsResponse, error) {
	node, err := s.factoryClient.Get(req.Id)
	if err != nil {
		log.Errorf("Failed to get node from factory: %v", err)
		return nil, err
	}

	if node.Node.OrgName == "" {
		log.Errorf("Node org name is empty")
		return nil, status.Errorf(codes.FailedPrecondition, "Node is not provisioned in any org")
	}
 
	dns := s.dnsMap[node.Node.OrgName]
	if dns == "" {
		log.Errorf("DNS is not found for org %s", node.Node.OrgName)
		return nil, status.Errorf(codes.NotFound, "DNS is not found for org %s", node.Node.OrgName)
	}	
 
	lookupSvc, err := s.lookupClient.GetClient()
	if err != nil {
		log.Errorf("Failed to get lookup client: %v", err)
		return nil, err
	}
	
	msgSystem, err := lookupSvc.GetSystemForOrg(ctx, &lpb.GetSystemRequest{
		OrgName:    node.Node.OrgName,
		SystemName: MessagingSystem,
	})
	if err != nil {
		log.Errorf("Failed to get messaging system: %v", err)
		return nil, err
	}
	
	ips, err := net.LookupIP(dns)
	if err != nil {
		log.Errorf("Could not get IPs: %v", err)
		return nil, err
	}
	
	var ip string
	for _, ipAddr := range ips {
		if ipv4 := ipAddr.To4(); ipv4 != nil {
			ip = ipv4.String()
			break
		}
	}

	if ip == "" {
		log.Errorf("No IPv4 address found for DNS %s", dns)
		return nil, status.Errorf(codes.NotFound, "No IPv4 address found for DNS %s", dns)
	}

	n, err := s.nodeRepo.GetNode(node.Node.Id)
	if err != nil && err.Error() != "record not found" {	
		log.Errorf("Failed to get node from database: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to get node from database: %v", err)
	}

	nd := &db.Node{}
	if n == nil {
		nd.NodeId = node.Node.Id
		nd.MeshPodName = ""
		nd.MeshPodIp = ""
		nd.MeshPodPort = 8082
	} else {
		nd.NodeId = n.NodeId
		nd.MeshPodName = n.MeshPodName
		nd.MeshPodIp = n.MeshPodIp
		nd.MeshPodPort = n.MeshPodPort
	}

	if err := utils.SpawnReplica(ctx, nd, s.config, s.clientSet, s.nodeRepo); err != nil {
		log.Warnf("Failed to spawn mesh replica for node %s: %v", nd.NodeId, err)
	}

	return &pb.GetNodeCredentialsResponse{
		Id:          node.Node.Id,
		OrgName:     node.Node.OrgName,
		Ip:          ip,
		Certificate: msgSystem.Certificate,
	}, nil
}

