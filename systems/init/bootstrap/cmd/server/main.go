/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package main

import (
	"os"

	"github.com/num30/config"

	"github.com/ukama/ukama/systems/init/bootstrap/cmd/version"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/server"
	"gopkg.in/yaml.v3"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/rest/client"
	factory "github.com/ukama/ukama/systems/common/rest/client/factory"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	provider "github.com/ukama/ukama/systems/init/bootstrap/client"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	runGrpcServer()
	log.Infof("Starting %s", pkg.ServiceName)
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if svcConf.DebugMode {
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = svcConf.DebugMode
}

func runGrpcServer() {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	factoryUrl, err := ic.GetHostUrl(ic.NewInitClient(svcConf.Http.InitClient, client.WithDebug()),
		ic.CreateHostString(svcConf.OrgName, "factory"), &svcConf.OrgName)
	if err != nil {
		log.Fatalf("Failed to resolve factory system address from initClient: %v", err)
	}
	factoryClient := factory.NewNodeFactoryClient(factoryUrl.String(), client.WithDebug())

	mbClient := mb.NewMsgBusClient(svcConf.MsgClient.Timeout, svcConf.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, svcConf.Queue.Uri, svcConf.Service.Uri, svcConf.MsgClient.Host, svcConf.MsgClient.Exchange, svcConf.MsgClient.ListenQueue, svcConf.MsgClient.PublishQueue, svcConf.MsgClient.RetryCount, svcConf.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	bootstrapServer := server.NewBootstrapServer(svcConf.OrgName, mbClient, svcConf.DebugMode, provider.NewLookupClientProvider(svcConf.Lookup, svcConf.Timeout), factoryClient)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterBootstrapServiceServer(s, bootstrapServer)

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

	log.Debug("awaiting terminate/interrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}
