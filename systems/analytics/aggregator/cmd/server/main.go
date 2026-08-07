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
	_ "time/tzdata" // embedded tzdata so org-timezone spans work on alpine

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/analytics/aggregator/cmd/version"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/performance"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/server"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
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

	// Aggregator owns migration of the rollup zone.
	err := d.Init(&schema.KpiRollup{})
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

	kpis, err := schema.LoadKpiSpecs(serviceConfig.Rollup.SpecsDir)
	if err != nil {
		log.Fatalf("Loading KPI specs from %s failed: %v", serviceConfig.Rollup.SpecsDir, err)
	}

	if len(kpis) == 0 {
		log.Fatalf("No KPI specs found in %s — is the configs dir present in the image?", serviceConfig.Rollup.SpecsDir)
	}

	log.Infof("Loaded %d KPI specs from %s", len(kpis), serviceConfig.Rollup.SpecsDir)

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName,
		pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	grid := schema.Grid{W: serviceConfig.Window.W}

	repo := db.NewRepo(sDb)

	engine, err := rollup.NewEngine(grid, kpis, repo, repo, serviceConfig.OrgName,
		serviceConfig.Rollup.Timezone, serviceConfig.Rollup.FlatThresholdPct,
		serviceConfig.Rollup.SweepInterval)
	if err != nil {
		log.Fatalf("Initializing rollup engine failed: %v", err)
	}

	reports, err := schema.LoadReportSpecs(serviceConfig.Rollup.ReportsDir, kpis)
	if err != nil {
		log.Fatalf("Loading report specs from %s failed: %v", serviceConfig.Rollup.ReportsDir, err)
	}

	log.Infof("Loaded %d report specs from %s", len(reports), serviceConfig.Rollup.ReportsDir)

	composer := performance.NewComposer(serviceConfig.OrgName, reports, repo, repo, repo, grid,
		serviceConfig.Rollup.ReportWindow)

	readServer := server.NewAggregatorServer(serviceConfig.OrgName, kpis, repo, composer, grid, repo)
	eventServer := server.NewAggregatorEventServer(serviceConfig.OrgName, engine)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		pb.RegisterAggregatorServiceServer(s, readServer)
		egenerated.RegisterEventNotificationServiceServer(s, eventServer)
	})

	grpcServer.RegisterDependency("db", true, ugrpc.DBCheck(sDb))
	grpcServer.RegisterDependency("msgclient", true, ugrpc.MsgClientCheck(serviceConfig.MsgClient.Host))

	go grpcServer.StartServer()
	go msgBusListener(mbClient)
	go engine.StartSweeper()

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
