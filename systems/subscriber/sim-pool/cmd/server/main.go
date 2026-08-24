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

	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	pkg "github.com/ukama/ukama/systems/subscriber/sim-pool/pkg"
)

const (
	factorySystem = "factory"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()
	simDb := initDb()
	runGrpcServer(simDb)
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

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Sim{})

	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
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

	factoryServiceUrl, err := ic.GetHostAddress(ic.NewInitClient(serviceConfig.Http.InitClient, cclient.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, factorySystem), &serviceConfig.OrgName)
	if err != nil {
		log.Fatalf("Failed to resolve %s system address from initClient: %v", factorySystem, err)
	}

	factoryClient := factory.NewSimFactoryClient(factoryServiceUrl.String(), cclient.WithDebug(serviceConfig.DebugMode))

	srv := server.NewSimPoolServer(serviceConfig.OrgName, db.NewSimRepo(gormdb), factoryClient, mbClient)
	nSrv := server.NewSimPoolEventServer(serviceConfig.OrgName, db.NewSimRepo(gormdb))

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		egenerated.RegisterEventNotificationServiceServer(s, nSrv)
		pb.RegisterSimServiceServer(s, srv)
	})

	grpcServer.RegisterDependency("db", true, ugrpc.DBCheck(gormdb))
	grpcServer.RegisterDependency("msgclient", true, ugrpc.MsgClientCheck(serviceConfig.MsgClient.Host))

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
