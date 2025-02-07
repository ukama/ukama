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

	"github.com/ukama/ukama/systems/data-plan/rate/pkg/client"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/server"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/data-plan/rate/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/db"
	"gopkg.in/yaml.v3"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mbc "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	generated "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"

	"google.golang.org/grpc"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	log.Infof("Starting the rate service")

	initConfig()
	rateDb := initDb()
	runGrpcServer(rateDb)
}

func initConfig() {

	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Markups{}, &db.DefaultMarkup{})

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

	mbClient := mbc.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	srv := server.NewRateServer(serviceConfig.OrgName, db.NewMarkupsRepo(gormdb), db.NewDefaultMarkupRepo(gormdb),
		client.NewBaseRateClientProvider(serviceConfig.BaseRate, serviceConfig.Timeout),
		mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterRateServiceServer(s, srv)
	})

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mbc.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
