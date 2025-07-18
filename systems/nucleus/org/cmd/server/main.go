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
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/nucleus/org/cmd/version"
	"github.com/ukama/ukama/systems/nucleus/org/pkg"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/db"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/providers"
	"github.com/ukama/ukama/systems/nucleus/org/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	orgDb := initDb()

	runGrpcServer(orgDb)
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if svcConf.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)

	err := d.Init(&db.Org{}, &db.User{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	orgDB := d.GetGormDb()

	initOrgDB(orgDB)

	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(svcConf.MsgClient.Timeout, svcConf.OrgName, pkg.SystemName,
		pkg.ServiceName, instanceId, svcConf.Queue.Uri, svcConf.Service.Uri, svcConf.MsgClient.Host,
		svcConf.MsgClient.Exchange, svcConf.MsgClient.ListenQueue, svcConf.MsgClient.PublishQueue,
		svcConf.MsgClient.RetryCount, svcConf.MsgClient.ListenerRoutes)

	user := providers.NewUserClientProvider(svcConf.UserHost)
	orch := providers.NewOrchestratorProvider(svcConf.OrchestratorHost, svcConf.DebugMode)
	registry := providers.NewRegistryProvider(svcConf.InitClientHost, svcConf.DebugMode)

	log.Debugf("MessageBus Client is %+v", mbClient)
	regServer := server.NewOrgServer(svcConf.OrgName, db.NewOrgRepo(gormdb),
		db.NewUserRepo(gormdb), orch, user, registry,
		mbClient, svcConf.Pushgateway, svcConf.DebugMode)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterOrgServiceServer(s, regServer)
	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	_ = regServer.PushMetrics()

	waitForExit()
}

func initOrgDB(orgDB *gorm.DB) {
	if orgDB.Migrator().HasTable(&db.Org{}) {
		if err := orgDB.First(&db.Org{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Initializing orgs table")

			var OwnerUUID, OrgUUID uuid.UUID
			var err error

			if OwnerUUID, err = uuid.FromString(svcConf.OwnerId); err != nil {
				log.Fatalf("Database initialization failed, need valid %v environment variable. Error: %v", "OWNERID", err)
			}

			if OrgUUID, err = uuid.FromString(svcConf.OrgId); err != nil {
				log.Fatalf("Database initialization failed, need valid %v environment variable. Error: %v", "ORGID", err)
			}

			org := &db.Org{
				Id:       OrgUUID,
				Owner:    OwnerUUID,
				Name:     svcConf.OrgName,
				Currency: svcConf.Currency,
				Country:  svcConf.Country,
			}

			usr := &db.User{
				Uuid: OwnerUUID,
			}

			if err := orgDB.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(usr).Error; err != nil {
					return err
				}

				u := &db.User{}
				if err := tx.First(&u, usr).Error; err != nil {
					return err
				}

				org.Users = append(org.Users, *u)
				if err := tx.Create(org).Error; err != nil {
					return err
				}

				// o := &db.Org{}
				// if err := tx.First(&o, org).Error; err != nil {
				// 	return err
				// }

				// usr.Org = []*db.Org{o}
				// if err := tx.Create(usr).Error; err != nil {
				// 	return err
				// }

				return nil
			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
		}
	}
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
