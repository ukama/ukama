/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package main

import (
	"os"

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/analytics/analysis/cmd/version"
	"github.com/ukama/ukama/systems/analytics/analysis/pkg"
	"github.com/ukama/ukama/systems/analytics/analysis/pkg/algos"
	"github.com/ukama/ukama/systems/analytics/analysis/pkg/db"
	"github.com/ukama/ukama/systems/analytics/analysis/pkg/engine"
	"github.com/ukama/ukama/systems/analytics/analysis/pkg/server"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	sDb := initDb()

	run(sDb)
}

func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)

	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err == nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	// Analysis owns migration of the KPI zone.
	err := d.Init(&schema.KpiWindow{}, &schema.AnalysisError{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func run(sDb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		instanceId = uuid.NewV4().String()
	}

	kpis, err := schema.LoadKpiSpecs(serviceConfig.Engine.SpecsDir)
	if err != nil {
		log.Fatalf("Loading KPI specs from %s failed: %v", serviceConfig.Engine.SpecsDir, err)
	}

	if len(kpis) == 0 {
		log.Fatalf("No KPI specs found in %s — is the configs dir present in the image?", serviceConfig.Engine.SpecsDir)
	}

	log.Infof("Loaded %d KPI specs from %s", len(kpis), serviceConfig.Engine.SpecsDir)

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName,
		pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	grid := schema.Grid{W: serviceConfig.Window.W}

	repo := db.NewRepo(sDb)

	runner, err := engine.NewRunner(grid, kpis, algos.Default(), repo, repo, repo, repo,
		mbClient, serviceConfig.OrgName, serviceConfig.Engine.SweepInterval,
		serviceConfig.Engine.CatchupWindows)
	if err != nil {
		log.Fatalf("Initializing analysis runner failed: %v", err)
	}

	eventServer := server.NewAnalysisEventServer(serviceConfig.OrgName, runner)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		egenerated.RegisterEventNotificationServiceServer(s, eventServer)
	})

	grpcServer.RegisterDependency("db", true, ugrpc.DBCheck(sDb))
	grpcServer.RegisterDependency("msgclient", true, ugrpc.MsgClientCheck(serviceConfig.MsgClient.Host))

	go grpcServer.StartServer()
	go msgBusListener(mbClient)
	go runner.StartSweeper()

	waitForExit()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start Message Client routine for service %s. Error %s",
			pkg.ServiceName, err.Error())
	}
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		log.Info(sig)
		done <- true
	}()

	<-done

	log.Infof("exiting service %s", pkg.ServiceName)
}
