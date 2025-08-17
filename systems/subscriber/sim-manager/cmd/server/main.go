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

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/interceptor"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
	cnuc "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	generated "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

const (
	registrySystemName     = "registry"
	dataplanSystemName     = "dataplan"
	notificationSystemName = "notification"
	ukamaAgentSystemName   = "ukamaagent"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", pkg.ServiceName)

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	simDB := initDb()

	runGrpcServer(simDB)

	log.Infof("Exiting service %s", pkg.ServiceName)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if serviceConfig.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("\nService: %s DB Config: %+v Service: %+v MsgClient Config %+v",
		pkg.ServiceName, serviceConfig.DB, serviceConfig.Service, serviceConfig.MsgClient)

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	err := d.Init(&db.Sim{}, &db.Package{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormDB sql.Db) {
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

	//TODO: We should do initclient resolution on demand, in order to avoid systems url changes side effects
	regUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, registrySystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	dataplanUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, dataplanSystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve dataplan address: %v", err)
	}

	notificationUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, notificationSystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve notification address: %v", err)
	}

	ukamaAgentUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, ukamaAgentSystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve ukama agent address: %v", err)
	}

	netClient := creg.NewNetworkClient(regUrl.String())
	pckgClient := cdplan.NewPackageClient(dataplanUrl.String())
	notificationClient := cnotif.NewMailerClient(notificationUrl.String())
	nucleusOrgClient := cnuc.NewOrgClient(serviceConfig.Http.NucleusClient)
	nucleusUserClient := cnuc.NewUserClient(serviceConfig.Http.NucleusClient)

	simManagerServer := server.NewSimManagerServer(
		serviceConfig.OrgName,
		db.NewSimRepo(gormDB),
		db.NewPackageRepo(gormDB),
		adapters.NewAgentFactory(serviceConfig.TestAgent, serviceConfig.OperatorAgent,
			ukamaAgentUrl.String(), serviceConfig.Timeout, pkg.IsDebugMode),
		pckgClient,
		providers.NewSubscriberRegistryClientProvider(serviceConfig.Registry, serviceConfig.Timeout),
		providers.NewSimPoolClientProvider(serviceConfig.SimPool, serviceConfig.Timeout),
		serviceConfig.Key,
		mbClient,
		serviceConfig.OrgId,
		serviceConfig.PushMetricHost,
		notificationClient,
		netClient,
		nucleusOrgClient,
		nucleusUserClient,
	)

	simManagerEventServer := server.NewSimManagerEventServer(serviceConfig.OrgName, simManagerServer)

	fsInterceptor := interceptor.NewFakeSimInterceptor(serviceConfig.TestAgent, serviceConfig.Timeout)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterSimManagerServiceServer(s, simManagerServer)
		egenerated.RegisterEventNotificationServiceServer(s, simManagerEventServer)
	})

	grpcServer.ExtraUnaryInterceptors = []grpc.UnaryServerInterceptor{
		fsInterceptor.UnaryServerInterceptor}

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
