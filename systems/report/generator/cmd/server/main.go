/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	_ "embed"
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/report/generator/cmd/version"
	"github.com/ukama/ukama/systems/report/generator/internal"
	"github.com/ukama/ukama/systems/report/generator/internal/pdf/engine"
	"github.com/ukama/ukama/systems/report/generator/internal/server"

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	fs "github.com/ukama/ukama/systems/report/generator/internal/server"
)

var serviceConfig = internal.NewConfig(internal.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(internal.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", internal.ServiceName)

	initConfig()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	runGrpcServer()

	log.Infof("Exiting service %s", internal.ServiceName)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	err := config.NewConfReader(internal.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if serviceConfig.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Debugf("\nService: %s DB Config: %+v Service: %+v MsgClient Config %+v",
		internal.ServiceName, serviceConfig.DB, serviceConfig.Service, serviceConfig.MsgClient)

	internal.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer() {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName,
		internal.SystemName, internal.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	pdfEngine, err := engine.NewWkGenerator()
	if err != nil {
		log.Fatalf("failed to get new PDF engine: %v", err)
	}

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		eSrv := server.NewGeneratorEventServer(serviceConfig.OrgName, pdfEngine, mbClient)
		egenerated.RegisterEventNotificationServiceServer(s, eSrv)
	})

	pdfServer := fs.NewPDFServer(serviceConfig.PdfHost, serviceConfig.PdfFolder,
		serviceConfig.PdfPrefix, serviceConfig.PdfPort)

	go pdfServer.Start()

	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s",
			internal.ServiceName, err.Error())
	}
}
