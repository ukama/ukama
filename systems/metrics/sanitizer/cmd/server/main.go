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
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/metrics/sanitizer/cmd/version"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	ic "github.com/ukama/ukama/systems/common/initclient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pkg "github.com/ukama/ukama/systems/metrics/sanitizer/pkg"
)

const (
	registrySystem = "registry"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	runGrpcServer()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	log.Infof("Initializing config")
	serviceConfig = pkg.NewConfig(pkg.ServiceName)

	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if pkg.IsDebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode

	log.Infof("Config: %+v", serviceConfig)
}

func runGrpcServer() {

	var mbClient mb.MsgBusServiceClient
	var instanceId string

	inst, ok := os.LookupEnv("POD_NAME")
	if ok {
		instanceId = inst
	} else {
		instanceId = pkg.InstanceId
	}

	if serviceConfig.IsMsgBus {
		mbClient = mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName,
			pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
			serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
			serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
			serviceConfig.MsgClient.RetryCount,
			serviceConfig.MsgClient.ListenerRoutes)

		log.Debugf("MessageBus Client is %+v", mbClient)
	}

	// Looking up registry system's host from initClient
	registrySystemUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, registrySystem),
		serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Fatalf("Failed to resolve %s system address from initClient: %v", registrySystem, err)
	}

	// Sanitizer service
	sanitizer, err := server.NewSanitizerServer(registrySystemUrl.String(), serviceConfig.PushGatewayHost,
		serviceConfig.OrgName, serviceConfig.Org, mbClient)
	if err != nil {
		log.Fatalf("Sanitizer server initialization failed. Error: %v", err)
	}

	rpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterSanitizerServiceServer(s, sanitizer)
	})

	if serviceConfig.IsMsgBus {
		go msgBusListener(mbClient)
	}

	rpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register with Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
