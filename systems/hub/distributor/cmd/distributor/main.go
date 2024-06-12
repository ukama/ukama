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
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/ukama/ukama/systems/common/config"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/hub/distributor/cmd/version"
	"github.com/ukama/ukama/systems/hub/distributor/pkg"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/distribution"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	generated "github.com/ukama/ukama/systems/hub/distributor/pb/gen"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/*Signal handler for SIGINT or SIGTERM to cancel a context in
	order to clean up and shut down gracefully if Ctrl+C is hit. */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* config parsig */
	initConfig()

	/* Log level */
	log.SetLevel(log.DebugLevel)

	/* Intilaize credentials */
	pkg.InitStoreCredentialsOptions(&serviceConfig.Storage)

	/* Start the HTTP server for chunk distribution */
	go startDistributionServer(ctx)

	/* Start the HTTP server for chunking request. */
	g := startChunkRequestServer()

	/* Signal Handling */
	handleSigterm(func() {
		log.Infof("Cleaning distribution service.")
		/* Call anything required for clean exit */
		g.StopServer()

		cancel()
	})

}

/* Start HTTP distribution server for distributing chunks */
func startDistributionServer(ctx context.Context) {
	err := distribution.RunDistribution(ctx, &serviceConfig.Distribution)
	if err != nil {
		log.Errorf("Error while starting distribution server : %s", err.Error())
		os.Exit(1)
	}
}

/* Start HTTP server for accepting chinking request from UkamaHub */
func startChunkRequestServer() *ugrpc.UkamaGrpcServer {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		instanceId = uuid.NewV4().String()
	}

	orgId, err := uuid.FromString(serviceConfig.OrgId)
	if err != nil {
		log.Fatalf("invalid org uuid. Error %s", err.Error())
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)
	log.Debugf("MessageBus Client is %+v", mbClient)

	distributorServer := server.NewDistributionServer(orgId, serviceConfig.OrgName, serviceConfig,
		mbClient, serviceConfig.PushGateway, serviceConfig.IsGlobal)

	log.Debugf("Distribution server is %+v and config %+v", distributorServer, serviceConfig.Grpc)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterDistributorServiceServer(s, distributorServer)
	})

	go grpcServer.StartServer()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	go msgBusListener(mbClient)

	return grpcServer
}

/* initConfig reads in config file, ENV variables, and flags if set. */
func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

/* Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting. */
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		handleExit()
		log.Infof("Exiting distribution service.")
		done <- true
	}()

	log.Debug("awaiting terminate/interrrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
