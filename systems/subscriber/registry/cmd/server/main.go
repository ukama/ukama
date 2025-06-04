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
	"os/signal"
	"syscall"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/subscriber/registry/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/client"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/registry/pkg/server"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	ic "github.com/ukama/ukama/systems/common/initclient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	registryDb := initDb()
	runGrpcServer(registryDb)
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
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}


func initDb() sql.Db {
    log.Infof("Initializing Database")
    d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
    err := d.Init(&db.Subscriber{})

    if err != nil {
        log.Fatalf("Database initialization failed. Error: %v", err)
    }
    
    if err := d.GetGormDb().Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_subscribers_active_email ON subscribers (email) WHERE deleted_at IS NULL").Error; err != nil {
        log.Fatalf("Failed to create conditional unique index on email. Error: %v", err)
    }
    
    return d
}
func setupSignalHandling(server *server.SubscriberServer) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        sig := <-c
        log.Infof("Received shutdown signal: %v", sig)
        
        server.Shutdown()
        
        log.Info("Graceful shutdown completed")
        os.Exit(0)
    }()
}


func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	regUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, "registry"), serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	networkClient := creg.NewNetworkClient(regUrl.String())

	orgClient := cnucl.NewOrgClient(serviceConfig.Http.NucleusClient)

	mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout,
		serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	simMClient := client.NewSimManagerClientProvider(serviceConfig.SimManagerHost)

	srv := server.NewSubscriberServer(serviceConfig.OrgName, db.NewSubscriberRepo(gormdb), mbClient, simMClient, serviceConfig.OrgId, orgClient, networkClient, serviceConfig)
	registryEventServer := server.NewRegistryEventServer(serviceConfig.OrgName, srv)
	
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		pb.RegisterRegistryServiceServer(s, srv)
			egenerated.RegisterEventNotificationServiceServer(s, registryEventServer)
	})
	 setupSignalHandling(srv)

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
