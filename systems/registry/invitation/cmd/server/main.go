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

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/invitation/cmd/version"
	"github.com/ukama/ukama/systems/registry/invitation/pkg"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/db"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
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

	notifUrl, err := ic.GetHostUrl(ic.CreateHostString(serviceConfig.OrgName, "notification"),
		serviceConfig.Http.InitClient, &serviceConfig.OrgName, serviceConfig.DebugMode)
	if err != nil {
		log.Fatalf("Failed to resolve notification system address from initClient: %v", err)
	}

	mailerClient := cnotif.NewMailerClient(notifUrl.String())
	orgClient := cnucl.NewOrgClient(serviceConfig.Http.NucleusClient)
	userClient := cnucl.NewUserClient(serviceConfig.Http.NucleusClient)

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName,
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

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register with Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
