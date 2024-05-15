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
	"github.com/ukama/ukama/systems/notification/distributor/cmd/version"
	"github.com/ukama/ukama/systems/notification/distributor/pkg"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	generated "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
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

	if pkg.IsDebugMode {
		log.SetLevel(log.DebugLevel)
	}
}

func runGrpcServer() {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	c := server.Clients{}

	c.Nucleus = providers.NewNucleusProvider(serviceConfig.Http.Nucleus, serviceConfig.DebugMode)

	c.Registry = providers.NewRegistryProvider(serviceConfig.Http.Nucleus, serviceConfig.DebugMode)

	c.Subscriber = providers.NewSubscriberProvider(serviceConfig.Http.Nucleus, serviceConfig.DebugMode)

	distributorServer := server.NewEventToNotifyServer(c, serviceConfig.OrgName, serviceConfig.OrgId, serviceConfig.DB, providers.NewEventNotifyClientProvider(serviceConfig.EventNotifyHost))

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterDistributorServiceServer(s, distributorServer)
	})

	go grpcServer.StartServer()

	waitForExit()
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
