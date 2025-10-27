/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"os"

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	ic "github.com/ukama/ukama/systems/common/initclient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	agent "github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/cmd/version"
	generated "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg/providers"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg/server"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	runGrpcServer()
}

func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer() {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout,
		serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	cdrC := clients.NewCDRClient(serviceConfig.Http.AgentNodeGateway)

	ukamaAgentUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, "ukamaagent"),
		serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Errorf("Failed to resolve ukama agent address: %v", err)
	}

	regUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, "registry"), serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	nodeClient := creg.NewNodeClient(regUrl.String())

	ukamaAgentClient := agent.NewUkamaAgentClient(ukamaAgentUrl.String())

	dsubServer := server.NewDsubscriberServer(serviceConfig.OrgName, mbClient, serviceConfig.RoutineConfig, nodeClient, cdrC, providers.NewDsimfactoryProvider(serviceConfig.DsimfactoryHost), ukamaAgentClient)
	nSrv := server.NewDsubEventServer(serviceConfig.OrgName, dsubServer)

	log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		egenerated.RegisterEventNotificationServiceServer(s, nSrv)
		generated.RegisterDsubscriberServiceServer(s, dsubServer)
	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	waitForExit()
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	go func() {

		sig := <-sigs
		log.Info(sig)
		done <- true
	}()

	log.Debug("awaiting terminate/interrrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}
