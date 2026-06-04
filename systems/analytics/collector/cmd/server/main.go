// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// Copyright (c) 2026-present, Ukama Inc.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/num30/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/analytics/collector/cmd/version"
	"github.com/ukama/ukama/systems/analytics/collector/pkg"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/db"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/refresh"
	"github.com/ukama/ukama/systems/analytics/collector/pkg/server"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
)

var serviceConfig *pkg.Config

func main() {
	initConfig()

	log.Infof("Starting %s version %s", pkg.ServiceName, version.Version)

	d := db.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	if err := d.Init(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	collectorRepo := db.NewCollectorRepo(d)
	snapshotRepo := db.NewSnapshotRepo(d)
	factRepo := db.NewFactRepo(d)
	stateRepo := db.NewRollupStateRepo(d)
	rollupRepo := db.NewRollupRepo(d)

	registryClient := refresh.NewRegistryClientSet(&serviceConfig.Http.RegistryClient)
	subscriberClient := refresh.NewSubscriberClientSet(&serviceConfig.Http.SubscriberClient)
	dataplanClient := refresh.NewDataplanClientSet(&serviceConfig.Http.DataplanClient)
	metricsClient := refresh.NewMetricsClientSet(&serviceConfig.Http.MetricsClient)
	nodeClient := refresh.NewNodeClientSet(&serviceConfig.Http.NodeClient)
	inventoryClient := refresh.NewInventoryClientSet(&serviceConfig.Http.InventoryClient)
	billingClient := refresh.NewBillingClientSet(&serviceConfig.Http.BillingClient)

	refresher := refresh.NewRefresher(refresh.RefresherConfig{
		OrgId:          serviceConfig.OrgId,
		OrgName:        serviceConfig.OrgName,
		RegistryClient: registryClient,
		SubscriberClient: subscriberClient,
		DataplanClient: dataplanClient,
		MetricsClient: metricsClient,
		NodeClient: nodeClient,
		InventoryClient: inventoryClient,
		BillingClient: billingClient,
		Snapshots: snapshotRepo,
		Facts: factRepo,
		Rollups: stateRepo,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := server.NewRollupScheduler(stateRepo, rollupRepo, server.RollupSchedulerConfig{
		Enabled:      serviceConfig.RollupScheduler.Enabled,
		Interval:     serviceConfig.RollupScheduler.Interval,
		LookbackDays: serviceConfig.RollupScheduler.LookbackDays,
	})
	scheduler.Start(ctx)

	runGrpcServer(collectorRepo, refresher, stateRepo, rollupRepo)

	if serviceConfig.Queue.Uri != "" {
		runMsgBusListener(collectorRepo, snapshotRepo, factRepo, stateRepo)
	} else {
		log.Warn("QUEUE_URI is empty; analytics event consumer disabled")
	}

	waitForExit()
}

func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)

	reader := config.NewConfReader(pkg.ServiceName)
	reader.SearchDirs = []string{"./config"}

	if err := reader.Read(serviceConfig); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if serviceConfig.DebugMode {
		log.SetLevel(log.DebugLevel)

		b, err := yaml.Marshal(serviceConfig)
		if err == nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
}

func runGrpcServer(
	collectorRepo db.CollectorRepo,
	refresher *refresh.Refresher,
	stateRepo db.RollupStateRepo,
	rollupRepo db.RollupRepo,
) {
	grpcServer := server.NewCollectorServer(
		collectorRepo,
		refresher,
		stateRepo,
		rollupRepo,
	)

	go grpcServer.StartServer()
}

func runMsgBusListener(
	collectorRepo db.CollectorRepo,
	snapshotRepo db.SnapshotRepo,
	factRepo db.FactRepo,
	stateRepo db.RollupStateRepo,
) {
	msgbusClient := msgbus.NewClient(serviceConfig.Queue.Uri)

	handler := server.NewEventHandler(server.EventHandlerConfig{
		OrgId:     serviceConfig.OrgId,
		OrgName:   serviceConfig.OrgName,
		Collector: collectorRepo,
		Snapshots:  snapshotRepo,
		Facts:      factRepo,
		Rollups:    stateRepo,
	})

	for _, consumer := range handler.Consumers() {
		if err := msgbusClient.Subscribe(
			consumer.Exchange,
			consumer.Queue,
			consumer.RoutingKey,
			consumer.Handler,
		); err != nil {
			log.Fatalf(
				"failed to subscribe exchange=%s queue=%s routingKey=%s: %v",
				consumer.Exchange,
				consumer.Queue,
				consumer.RoutingKey,
				err,
			)
		}

		log.Infof(
			"subscribed exchange=%s queue=%s routingKey=%s",
			consumer.Exchange,
			consumer.Queue,
			consumer.RoutingKey,
		)
	}

	msgbusClient.StartConsuming()
}

func waitForExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	go func() {
		sig := <-sigs
		log.Info(sig)
		done <- true
	}()

	log.Debug("awaiting terminate/interrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}

func init() {
	ukama.Init()
}
