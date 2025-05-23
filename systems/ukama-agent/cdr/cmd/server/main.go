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
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pb/gen"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/client"
	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/server"
	"gopkg.in/yaml.v3"

	pkg "github.com/ukama/ukama/systems/ukama-agent/cdr/pkg"

	"github.com/ukama/ukama/systems/ukama-agent/cdr/cmd/version"

	"github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egen "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {

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
			log.Infof("Config:\n %s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode

	if serviceConfig.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Config: %+v and routes %+v", serviceConfig, serviceConfig.MsgClient.ListenerRoutes)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, true)
	err := d.Init(&db.CDR{}, &db.Usage{})
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

	cdr := db.NewCDRRepo(gormdb)
	usage := db.NewUsageRepo(gormdb)

	// asr service
	asrClient, err := client.NewAsrClient(serviceConfig.AsrHost, serviceConfig.Timeout)
	if err != nil {
		log.Fatalf("ASR Client initilization failed. Error: %v", err)
	}

	cdrServer, err := server.NewCDRServer(cdr, usage, serviceConfig.OrgId, serviceConfig.OrgName,
		serviceConfig.PushGatewayHost, asrClient, mbClient)
	if err != nil {
		log.Fatalf("asr server initialization failed. Error: %v", err)
	}

	nSrv := server.NewCDREventServer(cdrServer, serviceConfig.OrgName)

	rpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterCDRServiceServer(s, cdrServer)
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
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
