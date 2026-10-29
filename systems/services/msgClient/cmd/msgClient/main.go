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
	"os/signal"
	"syscall"

	"github.com/num30/config"
	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/services/msgClient/cmd/version"
	"github.com/ukama/ukama/systems/services/msgClient/internal"
	"github.com/ukama/ukama/systems/services/msgClient/internal/broker"
	"github.com/ukama/ukama/systems/services/msgClient/internal/delivery"
	"github.com/ukama/ukama/systems/services/msgClient/internal/publish"
	"github.com/ukama/ukama/systems/services/msgClient/internal/registry"
	"github.com/ukama/ukama/systems/services/msgClient/internal/server"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	msgbus "github.com/ukama/ukama/systems/common/msgbus"
	generated "github.com/ukama/ukama/systems/common/pb/gen/msgclient"

	"google.golang.org/grpc"
)

var serviceConfig = internal.NewConfig()

func main() {
	ccmd.ProcessVersionArgument("msgClient", os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting the msgClient service")

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	runGrpcServer()

	log.Infof("Exiting service %s", internal.ServiceName)

}

func initConfig() {
	log.Infof("Initializing config")
	serviceConfig = &internal.Config{
		Grpc: &uconf.Grpc{
			Port: 9095,
		},
	}

	err := config.NewConfReader(internal.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if internal.IsDebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("Service: %s Config: %+v", internal.ServiceName, serviceConfig.Grpc)

}

func runGrpcServer() {
	b, err := broker.New(serviceConfig.Queue.Uri, serviceConfig.MsgBus.ManagementUri, serviceConfig.MsgBus.User, serviceConfig.MsgBus.Password)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ. Error: %s", err.Error())
	}

	pub, err := publish.New(b.Publish, serviceConfig.Publish.Timeout)
	if err != nil {
		log.Fatalf("Failed to create publisher. Error: %s", err.Error())
	}

	startListener := func(service, uri string, onUnreachable func()) (registry.Listener, error) {
		l, err := delivery.Start(b.Listen, service, uri, serviceConfig.Delivery.Timeout, onUnreachable)
		if err != nil {
			return nil, err
		}
		return l, nil
	}

	reg := registry.New(b, startListener, registry.Config{
		RefreshInterval: serviceConfig.Lease.RefreshInterval,
		MissedRefreshes: serviceConfig.Lease.MissedRefreshes,
		QueueMaxLength:  serviceConfig.MsgBus.QueueMaxLength,
		FlushAfter:      serviceConfig.MsgBus.QueueFlushAfter,
	})

	if err := reg.LoadOwned(); err != nil {
		log.Warnf("Failed to list existing service queues. Error: %s", err.Error())
	}

	stop := make(chan struct{})
	go reg.Run(stop)

	p := msgbus.NewShovelProvider(serviceConfig.MsgBus.ManagementUri, serviceConfig.DebugMode, serviceConfig.OrgName, serviceConfig.MsgBus.User, serviceConfig.MsgBus.Password,
		serviceConfig.Shovel.SrcUri, serviceConfig.Shovel.DestUri, serviceConfig.Shovel.DestExchange,
		serviceConfig.Shovel.SrcExchange, serviceConfig.Shovel.SrcExchangeKey)
	/* Create a shovel if required */
	initShovel(p)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewMsgClientServer(reg, pub, p, serviceConfig.System, serviceConfig.Lease.RefreshInterval)
		generated.RegisterMsgClientServiceServer(s, srv)
	})

	// grpcServer.RegisterDependency("rabbitmq", true, ugrpc.AmqpCheck(serviceConfig.Queue.Uri))

	signalHandler(reg, grpcServer, stop)

	grpcServer.StartServer()
}

func initShovel(p msgbus.MsgBusShovelProvider) {

	if serviceConfig.OrgName == serviceConfig.MasterOrgName {
		log.Infof("Master org %s running no need to add shovel.", serviceConfig.MasterOrgName)
		return
	}

	err := p.CreateShovel(serviceConfig.OrgName, nil)
	if err != nil {
		log.Fatalf("Failed to create shovelwith name %s. Error %+v.", serviceConfig.OrgName, err)
	}
}

func signalHandler(reg *registry.Registry, server *ugrpc.UkamaGrpcServer, stop chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		close(stop)
		reg.Close()
		server.StopServer()
	}()
}
