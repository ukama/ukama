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

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/notification/distributor/cmd/version"
	"github.com/ukama/ukama/systems/notification/distributor/pkg"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/db"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	sreg "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	generated "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
)

const (
	registrySystemName   = "registry"
	subscriberSystemName = "subscriber"
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
	log.Debugf("Distributor config %+v", serviceConfig)

	//TODO: we should do initclient resolution on demand, in order to avoid URL changes side effects.
	regUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, registrySystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	subUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, subscriberSystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	networkClient := creg.NewNetworkClient(regUrl.String())
	memberClient := creg.NewMemberClient(regUrl.String())
	subClient := sreg.NewSubscriberClient(subUrl.String())
	eventNotifyService := providers.NewEventNotifyClientProvider(serviceConfig.EventNotifyHost)

	nh := db.NewNotifyHandler(serviceConfig.DB, eventNotifyService)

	distributorServer := server.NewDistributorServer(networkClient, memberClient, subClient,
		nh, serviceConfig.OrgName, serviceConfig.OrgId, eventNotifyService)

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
