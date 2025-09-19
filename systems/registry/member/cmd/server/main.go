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

	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/member/cmd/version"
	"github.com/ukama/ukama/systems/registry/member/pkg"
	"github.com/ukama/ukama/systems/registry/member/pkg/db"
	"github.com/ukama/ukama/systems/registry/member/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	generated "github.com/ukama/ukama/systems/registry/member/pb/gen"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	mDb := initDb()
	runGrpcServer(mDb)
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
	err := d.Init(&db.Member{})
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

	id, err := uuid.FromString(serviceConfig.OrgId)
	if err != nil {
		log.Fatalf("invalid org uuid. Error %s", err.Error())
	}

	orgClient := cnucl.NewOrgClient(serviceConfig.Http.NucleusClient)
	userClient := cnucl.NewUserClient(serviceConfig.Http.NucleusClient)

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName,
		instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)
	memberServer := server.NewMemberServer(serviceConfig.OrgName, db.NewMemberRepo(gormdb),
		orgClient, userClient, mbClient, serviceConfig.PushGateway, id)

	memberEventServer := server.NewPackageEventServer(serviceConfig.OrgName, memberServer, serviceConfig.MasterOrgName)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterMemberServiceServer(s, memberServer)
		egenerated.RegisterEventNotificationServiceServer(s, memberEventServer)
	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	_ = memberServer.PushOrgMemberCountMetric(id)

	initMemberDB(gormdb, orgClient, userClient)

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

func initMemberDB(d sql.Db, orgClient cnucl.OrgClient, userClient cnucl.UserClient) {
	mDB := d.GetGormDb()
	if mDB.Migrator().HasTable(&db.Member{}) {
		if err := mDB.First(&db.Member{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Initializing registry member table for org")

			var OwnerUUID uuid.UUID
			var err error

			if OwnerUUID, err = uuid.FromString(serviceConfig.OwnerId); err != nil {
				log.Fatalf("Database initialization failed, need valid %v environment variable. Error: %v", "OWNERID", err)
			}

			user, err := userClient.GetById(serviceConfig.OwnerId)
			if err != nil {
				log.Fatalf("Failed to connect to user service for validation of owner %s. Error: %v", serviceConfig.OwnerId, err)
			}

			org, err := orgClient.Get(serviceConfig.OrgName)
			if err != nil {
				log.Fatalf("Failed to connect to org service for validation of owner %s. Error: %v", serviceConfig.OrgName, err)
			}

			if user.Id != org.Owner {
				log.Fatalf("Failed to validate user %s as owner of org %+v.", serviceConfig.OwnerId, org)
			}

			if user.IsDeactivated {
				log.Fatalf("User is %s is in %s state", serviceConfig.OwnerId, "deactivated")
			}

			if org.IsDeactivated {
				log.Fatalf("Org is %s in %s state", serviceConfig.OwnerId, "deactivated")
			}

			member := &db.Member{
				UserId:      OwnerUUID,
				Deactivated: false,
				MemberId:    uuid.NewV4(),
				Role:        roles.RoleType(roles.TYPE_OWNER),
			}
			if err := mDB.Transaction(func(tx *gorm.DB) error {
				err := orgClient.AddUser(org.Id, member.UserId.String())
				if err != nil {
					return err
				}
				if err := tx.Create(member).Error; err != nil {
					return err
				}
				return nil

			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
		}
	}
}
