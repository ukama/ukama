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
	"strings"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	factory "github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/init/bootstrap/client"
	"github.com/ukama/ukama/systems/init/bootstrap/utils"

	messaging "github.com/ukama/ukama/systems/common/rest/client/messaging"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)
 
const MessagingSystem = "messaging"

// BootstrapServerConfig holds required config and messaging cert for the server.
type BootstrapServerConfig struct {
	Config        *pkg.Config
	MessagingCert string
}

// BootstrapTestDeps holds optional deps for testing (nil in production).
// When ClientSet is non-nil, in-cluster k8s config is skipped. When DNSMap is non-nil, it overrides config.ToDNSMap().
type BootstrapTestDeps struct {
	ClientSet kubernetes.Interface
	DNSMap    map[string]string
}

type BootstrapServer struct {
	pb.UnimplementedBootstrapServiceServer
	bootstrapRoutingKey msgbus.RoutingKeyBuilder
	msgbus              mb.MsgBusServiceClient
	debug               bool
	lookupClient        client.LookupClientProvider
	factoryClient       factory.NodeFactoryClient
	dnsMap              map[string]string
	clientSet 			kubernetes.Interface
	config 				*pkg.Config
	nnsClient     		messaging.NnsClient
	messagingCert 		string
}
 
// NewBootstrapServerWithDeps creates a BootstrapServer; testDeps is nil in production.
func NewBootstrapServerWithDeps(msgBus mb.MsgBusServiceClient, debug bool, lookupClient client.LookupClientProvider, factoryClient factory.NodeFactoryClient, nnsClient messaging.NnsClient, serverConfig *BootstrapServerConfig, testDeps *BootstrapTestDeps) *BootstrapServer {
	var cs kubernetes.Interface
	var dns map[string]string
	if testDeps != nil && testDeps.ClientSet != nil {
		cs = testDeps.ClientSet
		if testDeps.DNSMap != nil {
			dns = testDeps.DNSMap
		} else {
			dns = serverConfig.Config.ToDNSMap()
		}
	} else {
		c, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		cs, err = kubernetes.NewForConfig(c)
		if err != nil {
			panic(err.Error())
		}
		dns = serverConfig.Config.ToDNSMap()
	}
	cfg := serverConfig.Config
	return &BootstrapServer{
		clientSet:           cs,
		bootstrapRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(cfg.OrgName).SetService(pkg.ServiceName),
		msgbus:              msgBus,
		debug:               debug,
		lookupClient:        lookupClient,
		factoryClient:       factoryClient,
		nnsClient:           nnsClient,
		dnsMap:              dns,
		config:              cfg,
		messagingCert:       serverConfig.MessagingCert,
	}
}

func NewBootstrapServer(msgBus mb.MsgBusServiceClient, debug bool, lookupClient client.LookupClientProvider, factoryClient factory.NodeFactoryClient, nnsClient messaging.NnsClient, config *pkg.Config, messagingCert string) *BootstrapServer {
	return NewBootstrapServerWithDeps(msgBus, debug, lookupClient, factoryClient, nnsClient, &BootstrapServerConfig{Config: config, MessagingCert: messagingCert}, nil)
}
 
func (s *BootstrapServer) GetNodeCredentials(ctx context.Context, req *pb.GetNodeCredentialsRequest) (*pb.GetNodeCredentialsResponse, error) {
	node, err := s.factoryClient.Get(req.Id)
	if err != nil {
		log.Errorf("Failed to get node from factory: %v", err)
		return nil, err
	}

	/* TODO: Uncomment this validation when init-lookup is ready for node lookup */

	// lookup, err := s.lookupClient.GetClient()
	// if err != nil {
	// 	log.Errorf("Failed to get node from lookup: %v", err)
	// 	return nil, err
	// }

	// lookupNode, err := lookup.GetNode(ctx, &lpb.GetNodeRequest{NodeId: req.Id})
	// if err != nil {
	// 	log.Errorf("Failed to get node from lookup: %v", err)
	// 	return nil, err
	// }

	if node.Node.OrgName == "" {
		log.Errorf("Node org name is empty")
		return nil, status.Errorf(codes.FailedPrecondition, "Node is not provisioned in any org")
	}
 
	dns := s.dnsMap[node.Node.OrgName]
	if dns == "" {
		log.Errorf("DNS is not found for org %s", node.Node.OrgName)
		return nil, status.Errorf(codes.NotFound, "DNS is not found for org %s", node.Node.OrgName)
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
 
	meshInfo, err := s.nnsClient.GetMesh(node.Node.Id)
	if err != nil && !strings.Contains(err.Error(), "node not found") {
		log.Errorf("Failed to get mesh info for node %s: %v", node.Node.Id, err)
		return nil, status.Errorf(codes.Internal, "Failed to get mesh info: %v", err)
	}
	
	nd :=utils.NodeMeshInfo{NodeId: node.Node.Id, MeshPodIp: "0.0.0.0", MeshPodPort: 8082}
	if meshInfo != nil && meshInfo.MeshIp != "" {
		log.Infof("Mesh info is found for node %s, mesh ip: %s", node.Node.Id, meshInfo.MeshIp)
		nd.MeshPodIp = meshInfo.MeshIp
		nd.MeshPodPort = int32(meshInfo.MeshPort)
	}

	if err := utils.SpawnReplica(ctx, nd, s.config, s.clientSet); err != nil {
		log.Warnf("Failed to spawn mesh replica for node %s: %v", node.Node.Id, err)
	}

	return &pb.GetNodeCredentialsResponse{
		Id:          node.Node.Id,
		OrgName:     node.Node.OrgName,
		Ip:          ip,
		Certificate: s.messagingCert,
	}, nil
}

