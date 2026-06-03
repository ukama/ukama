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

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/analytics/collector/cmd/version"
	"github.com/ukama/ukama/systems/analytics/collector/pkg"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/clients"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/refresh"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/server"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	generated "github.com/ukama/ukama/systems/analytics/collector/pb/gen"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	collectorDb := initDb()
	runGrpcServer(collectorDb)
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
		log.SetLevel(log.DebugLevel)
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)

	/* The collector owns the shared analytics schema: migrate ALL models. */
	err := d.Init(
		/* foundation */
		&db.EventLog{}, &db.EventError{}, &db.RefreshState{}, &db.RollupState{},
		/* snapshots */
		&db.NetworkSnapshot{}, &db.SiteSnapshot{}, &db.NodeSnapshot{},
		&db.CustomerSnapshot{}, &db.SimSnapshot{}, &db.SimBatchSnapshot{},
		&db.PackageSnapshot{}, &db.InventorySnapshot{}, &db.BillingSnapshot{},
		&db.HealthReportSnapshot{},
		/* facts */
		&db.PaymentEvent{}, &db.UsageEvent{}, &db.MetricSample{}, &db.AlarmEvent{},
		&db.NodeStateEvent{}, &db.SiteStateEvent{}, &db.CustomerEvent{},
		&db.SimEvent{}, &db.PackageEvent{}, &db.InventoryEvent{},
		/* intervals */
		&db.NodeStateInterval{}, &db.SiteStateInterval{},
		&db.CustomerPackageInterval{}, &db.SimStateInterval{},
		&db.MaintenanceWindow{},
		/* rollups */
		&db.BusinessSalesRollupDaily{}, &db.BusinessPackageRollupDaily{},
		&db.BusinessSiteRollupDaily{}, &db.BusinessInventoryRollupDaily{},
		&db.BusinessBillingRollupDaily{},
		&db.CustomerUsageRollupDaily{}, &db.CustomerStateRollupDaily{},
		&db.NetworkHealthRollupHourly{}, &db.SiteHealthRollupHourly{},
		&db.NodeHealthRollupHourly{}, &db.MetricRollupHourly{},
		&db.AlarmRollupDaily{}, &db.RadioRollupHourly{},
		&db.BackhaulRollupHourly{}, &db.PowerRollupHourly{},
	)
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	eventRepo := db.NewEventRepo(gormdb)
	stateRepo := db.NewStateRepo(gormdb)
	snapshotRepo := db.NewSnapshotRepo(gormdb)
	factRepo := db.NewFactRepo(gormdb)
	rollupRepo := db.NewRollupRepo(gormdb)

	refresher := refresh.NewRefresher(stateRepo, snapshotRepo, factRepo,
		clients.NewRegistryClient(serviceConfig.Http.RegistryClient),
		clients.NewSubscriberClient(serviceConfig.Http.SubscriberClient),
		clients.NewDataplanClient(serviceConfig.Http.DataplanClient),
		clients.NewMetricsClient(serviceConfig.Http.MetricsClient),
		clients.NewNodeClient(serviceConfig.Http.NodeClient),
		clients.NewInventoryClient(serviceConfig.Http.InventoryClient),
		clients.NewBillingClient(serviceConfig.Http.BillingClient))

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout,
		serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	collectorServer := server.NewCollectorServer(serviceConfig.OrgName, stateRepo,
		rollupRepo, eventRepo, refresher, mbClient, serviceConfig.PushGateway)

	collectorEventServer := server.NewCollectorEventServer(serviceConfig.OrgName,
		eventRepo, stateRepo, snapshotRepo, factRepo)

	scheduler := server.NewRollupScheduler(stateRepo, rollupRepo, server.RollupSchedulerConfig{
		Enabled:      serviceConfig.RollupScheduler.Enabled,
		Interval:     serviceConfig.RollupScheduler.Interval,
		LookbackDays: serviceConfig.RollupScheduler.LookbackDays,
	})
	scheduler.Start(context.Background())

	log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterCollectorServiceServer(s, collectorServer)
		egenerated.RegisterEventNotificationServiceServer(s, collectorEventServer)
	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	waitForExit()
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
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

	log.Debug("awaiting terminate/interrrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}
