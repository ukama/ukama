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
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/notify/cmd/version"
	"github.com/ukama/ukama/systems/notification/notify/internal"
	"github.com/ukama/ukama/systems/notification/notify/internal/db"
	"github.com/ukama/ukama/systems/notification/notify/internal/server"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	generated "github.com/ukama/ukama/systems/notification/notify/pb/gen"
)

var serviceConfig = internal.NewConfig(internal.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(internal.ServiceName, os.Args, version.Version)

	initConfig()
	metrics.StartMetricsServer(serviceConfig.Metrics)

	nodeDb := initDb()
	runGrpcServer(nodeDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(internal.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	internal.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	err := d.Init(&db.Notification{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")

	if instanceId == "" {
		/* used on local machines */
		instanceId = uuid.NewV4().String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, internal.SystemName,
		internal.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewNotifyServer(serviceConfig.OrgName, db.NewNotificationRepo(gormdb), mbClient)
		generated.RegisterNotifyServiceServer(s, srv)

		eSrv := server.NewNotifyEventServer(serviceConfig.OrgName, db.NewNotificationRepo(gormdb), mbClient,
			serviceConfig.MsgClient.ListenerRoutes)
		egenerated.RegisterEventNotificationServiceServer(s, eSrv)

	})

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s",
			internal.ServiceName, err.Error())
	}
}
