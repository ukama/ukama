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

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/invitation/cmd/version"
	"github.com/ukama/ukama/systems/registry/invitation/pkg"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/db"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	generated "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	invitationDb := initDb()
	runGrpcServer(invitationDb)
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
	err := d.Init(&db.Invitation{})
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

	mailerClient := cnotif.NewMailerClient(serviceConfig.NotificationHost)
	orgClient := cnucl.NewOrgClient(serviceConfig.OrgRegistryHost)
	userClient := cnucl.NewUserClient(serviceConfig.OrgRegistryHost)

	mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName,
		pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri,
		serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange, serviceConfig.MsgClient.ListenQueue,
		serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	invitationServer := server.NewInvitationServer(db.NewInvitationRepo(gormdb),
		serviceConfig.InvitationExpiryTime, serviceConfig.AuthLoginbaseURL,
		mailerClient, orgClient, userClient, mbClient, serviceConfig.OrgName, serviceConfig.TemplateName)

	log.Debugf("MessageBus Client is %+v", mbClient)
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterInvitationServiceServer(s, invitationServer)
	})

	go grpcServer.StartServer()

	waitForExit()
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
