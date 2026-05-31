/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/num30/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"github.com/ukama/ukama/systems/operation/manager/pkg"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
	"github.com/ukama/ukama/systems/operation/manager/pkg/server"

	"github.com/ukama/ukama/systems/operation/manager/cmd/version"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	if pkg.InstanceId == "" {
		pkg.InstanceId = uuid.NewV4().String()
	}

	initConfig()
	gormDb := initDb()
	runGrpcServer(gormDb)
	log.Infof("Starting %s/%s", pkg.SystemName, pkg.ServiceName)
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	if err := config.NewConfReader(pkg.ServiceName).Read(svcConf); err != nil {
		log.Fatal("Error reading config ", err)
	} else if svcConf.DebugMode {
		if b, err := yaml.Marshal(svcConf); err == nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	if err := d.Init(&db.Operation{}, &db.ResourceLock{}, &db.OperationAudit{}); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	return d
}

func runGrpcServer(gormDb sql.Db) {
	mbClient := mb.NewMsgBusClient(
		svcConf.MsgClient.Timeout, svcConf.OrgName, pkg.SystemName, pkg.ServiceName,
		pkg.InstanceId, svcConf.Queue.Uri, svcConf.Service.Uri, svcConf.MsgClient.Host,
		svcConf.MsgClient.Exchange, svcConf.MsgClient.ListenQueue,
		svcConf.MsgClient.PublishQueue, svcConf.MsgClient.RetryCount,
		svcConf.MsgClient.ListenerRoutes,
	)
	log.Debugf("MessageBus Client is %+v", mbClient)

	repo := db.NewOperationRepo(gormDb)
	opServer := server.NewOperationServer(svcConf.OrgName, svcConf.OrgId, repo, mbClient)
	eventServer := server.NewEventServer(svcConf.OrgName, repo)

	grpcSrv := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterOperationManagerServiceServer(s, opServer)
		egenerated.RegisterEventNotificationServiceServer(s, eventServer)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.NewSweeper(repo).Run(ctx)

	go grpcSrv.StartServer()
	go msgBusListener(mbClient)

	waitForExit()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service: %v", err)
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start MsgBus listener: %v", err)
	}
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	log.Info("received signal: ", sig)
	log.Infof("exiting service %s/%s", pkg.SystemName, pkg.ServiceName)
}
