/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"errors"
	"os"

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/configurator/cmd/version"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"
	"github.com/ukama/ukama/systems/node/configurator/pkg/providers"
	"github.com/ukama/ukama/systems/node/configurator/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/rest/client"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	configstore "github.com/ukama/ukama/systems/node/configurator/pkg/configStore"
)

const registrySystemName = "registry"

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	mDb := initDb()
	log.Infof("Starting %s", pkg.ServiceName)
	runGrpcServer(mDb)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Configuration{}, &db.Commit{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	log.Infof("serviceConfig %+v", serviceConfig)
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

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	//TODO: We should do initclient resolution on demand, in order to avoid systems url changes side effects
	regUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, registrySystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	cnet := creg.NewNetworkClient(regUrl.String())
	csite := creg.NewSiteClient(regUrl.String())
	cnode := creg.NewNodeClient(regUrl.String())

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName,
		pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri,
		serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange, serviceConfig.MsgClient.ListenQueue,
		serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	s, err := providers.NewStoreClient(serviceConfig.StoreUrl, serviceConfig.StoreUser, serviceConfig.AccessToken, serviceConfig.Timeout)
	if err != nil {
		log.Fatalf("Failed to create a config store client. Error %s", err.Error())
	}
	configStore := configstore.NewConfigStore(mbClient, cnet, csite, cnode, db.NewConfigRepo(gormdb),
		db.NewCommitRepo(gormdb), serviceConfig.OrgName, s, serviceConfig.Timeout)

	configuratorServer := server.NewConfiguratorServer(mbClient, db.NewConfigRepo(gormdb), db.NewCommitRepo(gormdb), configStore,
		serviceConfig.OrgName, pkg.IsDebugMode)

	configuratorEventServer := server.NewConfiguratorEventServer(serviceConfig.OrgName, configuratorServer)

	log.Debugf("MessageBus Client config: %+v Client: %+v", serviceConfig.MsgClient, mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		pb.RegisterConfiguratorServiceServer(s, configuratorServer)
		epb.RegisterEventNotificationServiceServer(s, configuratorEventServer)
	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	initCommitDB(gormdb)

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

func initCommitDB(d sql.Db) {
	mDB := d.GetGormDb()
	if mDB.Migrator().HasTable(&db.Commit{}) {
		if err := mDB.First(&db.Commit{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Initializing commit table for configurator")

			/* TODO: validate the Hash */
			commit := &db.Commit{
				Hash: serviceConfig.LatestConfigHash,
			}

			if err := mDB.Transaction(func(tx *gorm.DB) error {

				if err := tx.Create(commit).Error; err != nil {
					return err
				}
				return nil

			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
		}
	}
}
