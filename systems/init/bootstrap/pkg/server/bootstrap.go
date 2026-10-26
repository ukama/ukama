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
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/mesh"

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
	mesh                *mesh.Manager
	config              *pkg.Config
	messagingCert       string
}

// NewBootstrapServerWithDeps creates a BootstrapServer; testDeps is nil in production.
func NewBootstrapServerWithDeps(msgBus mb.MsgBusServiceClient, debug bool, lookupClient client.LookupClientProvider, factoryClient factory.NodeFactoryClient, serverConfig *BootstrapServerConfig, testDeps *BootstrapTestDeps) *BootstrapServer {
	cfg := serverConfig.Config
	var cs kubernetes.Interface
	var dns map[string]string
	if testDeps != nil && testDeps.ClientSet != nil {
		cs = testDeps.ClientSet
		if testDeps.DNSMap != nil {
			dns = testDeps.DNSMap
		} else {
			dns = cfg.ToDNSMap()
		}
	} else {
		c, err := rest.InClusterConfig()
		if err != nil {
			log.Warnf("Not running in Kubernetes (in-cluster config unavailable: %v). Mesh pod spawning will be disabled.", err)
			cs = nil
		} else {
			// client-go defaults to 5 QPS, which a burst of booting nodes exhausts.
			if cfg.Mesh.KubeQps > 0 {
				c.QPS = float32(cfg.Mesh.KubeQps)
			}
			if cfg.Mesh.KubeBurst > 0 {
				c.Burst = cfg.Mesh.KubeBurst
			}
			cs, err = kubernetes.NewForConfig(c)
			if err != nil {
				log.Fatalf("Failed to create Kubernetes client: %v", err)
			}
		}
		dns = cfg.ToDNSMap()
	}

	var meshManager *mesh.Manager
	if cs != nil {
		meshManager = mesh.NewManager(cs, mesh.Config{
			Namespace:       cfg.OrgName + "-" + cfg.MeshNamespace,
			OrgName:         cfg.OrgName,
			TemplateRefresh: cfg.Mesh.TemplateRefresh,
		})
	}

	return &BootstrapServer{
		mesh:                meshManager,
		bootstrapRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(cfg.OrgName).SetService(pkg.ServiceName),
		msgbus:              msgBus,
		debug:               debug,
		lookupClient:        lookupClient,
		factoryClient:       factoryClient,
		dnsMap:              dns,
		config:              cfg,
		messagingCert:       serverConfig.MessagingCert,
	}
}

func NewBootstrapServer(msgBus mb.MsgBusServiceClient, debug bool, lookupClient client.LookupClientProvider, factoryClient factory.NodeFactoryClient, config *pkg.Config, messagingCert string) *BootstrapServer {
	return NewBootstrapServerWithDeps(msgBus, debug, lookupClient, factoryClient, &BootstrapServerConfig{Config: config, MessagingCert: messagingCert}, nil)
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

	if s.mesh != nil {
		if err := s.mesh.Ensure(ctx, node.Node.Id); err != nil {
			log.Warnf("Mesh for node %s is not ready: %v", node.Node.Id, err)
			return nil, err
		}
	}

	return &pb.GetNodeCredentialsResponse{
		Id:          node.Node.Id,
		OrgName:     node.Node.OrgName,
		Ip:          ip,
		Certificate: s.messagingCert,
	}, nil
}
