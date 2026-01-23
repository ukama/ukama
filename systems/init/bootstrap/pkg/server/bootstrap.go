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

	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/db"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	pod, err := s.spawnReplica(ctx, node.Node.Id)
	if err != nil {
		log.Errorf("Failed to spawn replica: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to spawn mesh pod replica: %v", err)
	}

	meshPodName := pod.Name

	err = s.nodeRepo.CreateNode(&db.Node{
		NodeId: node.Node.Id,
		MeshPodName: meshPodName,
	})
	if err != nil {
		log.Errorf("Failed to create node: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to create node: %v", err)
	}

	return &pb.GetNodeCredentialsResponse{
		Id:          node.Node.Id,
		OrgName:     node.Node.OrgName,
		Ip:          ip,
		Certificate: msgSystem.Certificate,
	}, nil
}

func (s *BootstrapServer) spawnReplica(ctx context.Context, nodeId string) (*corev1.Pod, error) {
	namespace := s.config.OrgName + "-" + s.config.MeshNamespace
	deployments, err := s.clientSet.AppsV1().Deployments(s.config.MeshNamespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: "app=mesh",
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not list deployments: %v", err)
	}

	if len(deployments.Items) == 0 {
		return nil, status.Errorf(codes.NotFound, "No deployment found with label app=mesh in namespace %s", s.config.MeshNamespace)
	}

	templateDeployment := &deployments.Items[0]
	podSpec := templateDeployment.Spec.Template.Spec.DeepCopy()

	podSpec.RestartPolicy = corev1.RestartPolicyOnFailure

	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "mesh-node-" + nodeId + "-",
			Namespace:    namespace,
		},
		Spec: *podSpec,
	}

	createdPod, err := s.clientSet.CoreV1().Pods(namespace).Create(
		context.TODO(),
		newPod,
		metav1.CreateOptions{},
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create pod: %v", err)
	}

	return createdPod, nil
}