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

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/notification/event-notify/cmd/version"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	sreg "github.com/ukama/ukama/systems/common/rest/client/subscriber"
	generated "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

const (
	registrySystemName   = "registry"
	subscriberSystemName = "subscriber"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	eventToNotifyDb := initDb()
	runGrpcServer(eventToNotifyDb)
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
	err := d.Init(&db.Notification{}, &db.Users{}, &db.UserNotification{}, &db.EventMsg{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	initTrigger(d)

	return d
}

func initTrigger(db sql.Db) {
	db.GetGormDb().Exec("CREATE OR REPLACE FUNCTION public.user_notifications_trigger() RETURNS TRIGGER AS $$ DECLARE notification_data text; BEGIN notification_data := NEW.id::text || ',' || NEW.notification_id::text || ',' || NEW.user_id::text || ',' || NEW.is_read::text; PERFORM pg_notify('user_notifications_channel', notification_data); RETURN NEW; END; $$ LANGUAGE plpgsql;")
	db.GetGormDb().Exec("DROP TRIGGER notify_trigger ON user_notifications;")
	db.GetGormDb().Exec("CREATE TRIGGER notify_trigger AFTER INSERT OR UPDATE ON user_notifications FOR EACH ROW EXECUTE FUNCTION public.user_notifications_trigger();")
}

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	//TODO:: need to do initclient resolutions on demand, in order to avoid url changes.
	regUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, registrySystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	subUrl, err := ic.GetHostUrl(ic.NewInitClient(serviceConfig.Http.InitClient, client.WithDebug(serviceConfig.DebugMode)),
		ic.CreateHostString(serviceConfig.OrgName, subscriberSystemName), &serviceConfig.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve subscriber address: %v", err)
	}

	orgClient := cnucl.NewOrgClient(serviceConfig.Http.NucleusClient)
	userClient := cnucl.NewUserClient(serviceConfig.Http.NucleusClient)
	memberClient := creg.NewMemberClient(regUrl.String())
	subscriberClient := sreg.NewSubscriberClient(subUrl.String())

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout,
		serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	eventToNotifyServer := server.NewEventToNotifyServer(serviceConfig.OrgName, serviceConfig.OrgId, memberClient, db.NewNotificationRepo(gormdb),
		db.NewUserRepo(gormdb), db.NewEventMsgRepo(gormdb), db.NewUserNotificationRepo(gormdb), mbClient)

	eventToNotifyEventServer := server.NewNotificationEventServer(serviceConfig.OrgName, serviceConfig.OrgId, subscriberClient, eventToNotifyServer)
	log.Debugf("MessageBus Client is %+v and config %+v", mbClient, serviceConfig.MsgClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterEventToNotifyServiceServer(s, eventToNotifyServer)
		egenerated.RegisterEventNotificationServiceServer(s, eventToNotifyEventServer)
	})

	go msgBusListener(mbClient)

	go grpcServer.StartServer()

	initUserDB(gormdb, orgClient, userClient)

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

func initUserDB(d sql.Db, orgClient cnucl.OrgClient, userClient cnucl.UserClient) {
	mDB := d.GetGormDb()
	if mDB.Migrator().HasTable(&db.Users{}) {
		if err := mDB.First(&db.Users{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Initializing user database for notification")

			var OwnerUUID uuid.UUID
			var err error

			if OwnerUUID, err = uuid.FromString(serviceConfig.OwnerId); err != nil {
				log.Fatalf("Database initialization failed, need valid %v environment variable. Error: %v", "OWNERID", err)
			}

			/* TODO: validate the user from user services */
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

			u := &db.Users{
				Id:     uuid.NewV4(),
				OrgId:  serviceConfig.OrgId,
				UserId: OwnerUUID.String(),
				Role:   roles.RoleType(roles.TYPE_OWNER),
			}
			if err := mDB.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(u).Error; err != nil {
					return err
				}
				return nil

			}); err != nil {
				log.Fatalf("Database initialization failed, invalid initial state. Error: %v", err)
			}
		}
	}
}
