/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/num30/config"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/analytics/analytics/cmd/version"
	"github.com/ukama/ukama/systems/analytics/analytics/pkg"
	businesspb "github.com/ukama/ukama/systems/analytics/business/pb/gen"
	businessdb "github.com/ukama/ukama/systems/analytics/business/pkg/db"
	businessserver "github.com/ukama/ukama/systems/analytics/business/pkg/server"
	customerpb "github.com/ukama/ukama/systems/analytics/customer/pb/gen"
	customerdb "github.com/ukama/ukama/systems/analytics/customer/pkg/db"
	customerserver "github.com/ukama/ukama/systems/analytics/customer/pkg/server"
	networkpb "github.com/ukama/ukama/systems/analytics/network/pb/gen"
	networkdb "github.com/ukama/ukama/systems/analytics/network/pkg/db"
	networkserver "github.com/ukama/ukama/systems/analytics/network/pkg/server"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	analyticsDb := initDb()
	runGrpcServer(analyticsDb)
}

func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	}

	if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err == nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing analytics read database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	// Read-only service: only connect, never AutoMigrate.
	// The collector service owns the analytics schema.
	if err := d.Connect(); err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormdb sql.Db) {
	businessServer := businessserver.NewBusinessServer(
		serviceConfig.OrgName,
		businessdb.NewSalesRepo(gormdb),
		businessdb.NewPackageRepo(gormdb),
		businessdb.NewSiteRepo(gormdb),
		businessdb.NewBillingRepo(gormdb),
		businessdb.NewInventoryRepo(gormdb),
		businessdb.NewActivityRepo(gormdb),
		nil,
		serviceConfig.PushGateway,
		serviceConfig.OrgId,
	)

	customerServer := customerserver.NewCustomerServer(
		serviceConfig.OrgName,
		customerdb.NewCustomerRepo(gormdb),
		customerdb.NewSimRepo(gormdb),
		customerdb.NewSupportRepo(gormdb),
		nil,
		serviceConfig.SimLowStockThreshold,
	)

	networkServer := networkserver.NewNetworkServer(
		serviceConfig.OrgName,
		networkdb.NewSiteRepo(gormdb),
		networkdb.NewNodeRepo(gormdb),
		networkdb.NewAlarmRepo(gormdb),
		networkdb.NewMetricRepo(gormdb),
		networkdb.NewEventRepo(gormdb),
		networkdb.NewHealthRepo(gormdb),
		nil,
		serviceConfig.PushGateway,
		serviceConfig.OrgId,
		serviceConfig.NetworkLatencyThresholdMs,
		serviceConfig.BatteryCriticalPercent,
		serviceConfig.TelemetryFreshSeconds,
	)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		businesspb.RegisterBusinessServiceServer(s, businessServer)
		customerpb.RegisterCustomerServiceServer(s, customerServer)
		networkpb.RegisterNetworkServiceServer(s, networkServer)
	})

	go grpcServer.StartServer()
	waitForExit()
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Infof("received signal %s", sig.String())
	log.Infof("exiting service %s", pkg.ServiceName)
}
