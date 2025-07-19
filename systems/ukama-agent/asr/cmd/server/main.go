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

	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/ukama-agent/asr/cmd/version"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egen "github.com/ukama/ukama/systems/common/pb/gen/events"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	pkg "github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

const (
	registrySystem = "registry"
	dataPlanSystem = "dataplan"
)

var serviceConfig *pkg.Config

func main() {
	log.SetLevel(log.DebugLevel)
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	hssDb := initDb()
	runGrpcServer(hssDb)
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

	if serviceConfig.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Config: %+v", serviceConfig)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, true)
	err := d.Init(&db.Asr{}, &db.Guti{}, &db.Tai{}, &db.Policy{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
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
	} else {
		log.Fatalf("MsgBus is mandatory for service %s", pkg.ServiceName)
	}

	asrRepo := db.NewAsrRecordRepo(gormdb)
	gutiRepo := db.NewGutiRepo(gormdb)
	//policyRepo := db.NewPolicyRepo(gormdb)

	// For now, we either assuming factory is global and/or currently using a dummy stub unter ukama/testing
	factory, err := client.NewFactoryClient(serviceConfig.FactoryHost, pkg.IsDebugMode)
	if err != nil {
		log.Fatalf("Factory Client initialization failed. Error: %v", err)
	}

	// Looking up registry system's host from initClient
	networkServiceUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, registrySystem),
		serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Fatalf("Failed to resolve %s system address from initClient: %v", registrySystem, err)
	}

	networkClient := registry.NewNetworkClient(networkServiceUrl.String(), cclient.WithDebug())

	// Looking up data plan system's host from initClient
	dataPlanUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, dataPlanSystem),
		serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Fatalf("Failed to resolve %s system address from initClient: %v", dataPlanSystem, err)
	}

	cdr, err := client.NewCDR(serviceConfig.CDRHost, serviceConfig.Timeout)
	if err != nil {
		log.Fatalf("CDR Client initilization failed. Error: %v", err)
	}

	//pcrf := pcrf.NewPCRFController(policyRepo, serviceConfig.DataplanHost, mbClient, serviceConfig.OrgName, serviceConfig.Reroute)

	controller := pm.NewPolicyController(asrRepo, mbClient, dataPlanUrl.String(),
		serviceConfig.OrgName, serviceConfig.OrgId, serviceConfig.Reroute, serviceConfig.Period, serviceConfig.Monitor)

	// ASR service
	asrServer, err := server.NewAsrRecordServer(asrRepo, gutiRepo,
		factory, networkClient, controller, cdr, serviceConfig.OrgId, serviceConfig.OrgName,
		mbClient, serviceConfig.AllowedTimeOfService) //
	if err != nil {
		log.Fatalf("asr server initialization failed. Error: %v", err)
	}

	nSrv := server.NewAsrEventServer(asrRepo, asrServer, gutiRepo, serviceConfig.OrgName)

	rpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterAsrRecordServiceServer(s, asrServer)
		if serviceConfig.IsMsgBus {
			egen.RegisterEventNotificationServiceServer(s, nSrv)
		}
	})

	if serviceConfig.IsMsgBus {
		go msgBusListener(mbClient)
	}

	rpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s",
			pkg.ServiceName, err.Error())
	}
}
