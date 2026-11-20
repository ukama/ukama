/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"context"
	"os"
	"strings"

	"github.com/num30/config"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/messaging/nns/cmd/version"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
	"github.com/ukama/ukama/systems/messaging/nns/pkg/server"

	dnspb "github.com/coredns/coredns/pb"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	nnsClient := pkg.NewNns(serviceConfig)

	metrics.StartMetricsServer(serviceConfig.Metrics)

	runGrpcServer(nnsClient)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	log.Infof("Initializing config")

	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		log.SetLevel(log.TraceLevel)
		log.Infof("Config is %+v DnsConfig %+v", serviceConfig, serviceConfig.Dns)
		// b, err := yaml.Marshal(serviceConfig)
		// if err != nil {
		// 	log.Infof("Config:\n%s", string(b))
		// }
	}

	log.Debugf("\nService: %s Service: %+v MsgClient Config %+v", pkg.ServiceName, serviceConfig.Service, serviceConfig.MsgClient)

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer(nns *pkg.Nns) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	regUrl, err := ic.GetHostAddress(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, "registry"), &serviceConfig.OrgName)
	if err != nil || regUrl == nil || regUrl.String() == "" {
		log.Fatalf("Failed to resolve registry address: %v", err)
	}

	nodeClient := creg.NewNodeClient(regUrl.String())

	srv := server.NewNnsServer(nns, serviceConfig, serviceConfig.Dns)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		eSrv := server.NewNnsEventServer(serviceConfig.OrgName, nodeClient, srv, serviceConfig.OrgId)
		pb.RegisterNnsServer(s, srv)
		dnspb.RegisterDnsServiceServer(s, server.NewDnsServer(nns, serviceConfig.Dns))
		egenerated.RegisterEventNotificationServiceServer(s, eSrv)
	})

	grpcServer.RegisterDependency("msgclient", true, ugrpc.MsgClientCheck(serviceConfig.MsgClient.Host))

	go msgBusListener(mbClient)

	startMeshWatcher(srv, mbClient)

	grpcServer.StartServer()
}

const namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func startMeshWatcher(srv *server.NnsServer, m mb.MsgBusServiceClient) {
	if !serviceConfig.MeshWatch {
		log.Infof("Mesh pod watch is disabled")
		return
	}

	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Warnf("Mesh pod watch not started: %v", err)
		return
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Warnf("Mesh pod watch not started: %v", err)
		return
	}

	namespace := serviceConfig.MeshNamespace
	if namespace == "" {
		b, err := os.ReadFile(namespaceFile)
		if err != nil {
			log.Warnf("Mesh pod watch not started: %v", err)
			return
		}
		namespace = strings.TrimSpace(string(b))
	}

	err = server.NewMeshWatcher(srv, m, serviceConfig.OrgName).Start(context.Background(), clientSet, namespace)
	if err != nil {
		log.Warnf("Mesh pod watch not started: %v", err)
	}
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
