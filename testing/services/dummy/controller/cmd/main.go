/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/num30/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"

	"github.com/ukama/ukama/testing/services/dummy/controller/cmd/version"

	generated "github.com/ukama/ukama/testing/services/dummy/controller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/controller/pkg"
	"github.com/ukama/ukama/testing/services/dummy/controller/pkg/server"
)
 
 var serviceConfig *pkg.Config
 
 func main() {
	 ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	 initConfig()
	 runGrpcServer()
 }
 
 func initConfig() {
	 serviceConfig = pkg.NewConfig(pkg.ServiceName)
	 err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	 if err != nil {
		 log.Fatalf("Error reading config: %v", err)
	 }
	 
	 if serviceConfig.DebugMode {
		 b, err := yaml.Marshal(serviceConfig)
		 if err != nil {
			 log.Warnf("Error marshaling config: %v", err)
		 } else {
			 log.Infof("Config:\n%s", string(b))
		 }
	 }
	 pkg.IsDebugMode = serviceConfig.DebugMode
 }
 
 func runGrpcServer() {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		inst := uuid.NewV4()
		instanceId = inst.String()
	}


	nodeClient := creg.NewNodeClient(serviceConfig.RegistryClient)


	
	 mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout,
		serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)
		controllerServer := server.NewControllerServer(serviceConfig.OrgName, serviceConfig.OrgId,nodeClient,serviceConfig.DnodeHost,mbClient)
		nSrv := server.NewEventServer(serviceConfig.OrgName, controllerServer)
   
	log.Debugf("MessageBus Client is %+v", mbClient)
	 grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		egenerated.RegisterEventNotificationServiceServer(s, nSrv)
		 generated.RegisterMetricsControllerServer(s, controllerServer)
	 })
 	go msgBusListener(mbClient)

	 go grpcServer.StartServer()
 
	 go startMetricsServer()
	 
	 sigs := make(chan os.Signal, 1)
	 signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	 
	 sig := <-sigs
	 log.Infof("Received signal %v, shutting down...", sig)
	 
	 controllerServer.Cleanup()
	 
	 log.Infof("Exiting service %s", pkg.ServiceName)
 }
 
 func startMetricsServer() {
	 mux := http.NewServeMux()
	 
	 mux.Handle("/metrics", promhttp.Handler())
	 
	 mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		 w.WriteHeader(http.StatusOK)
		 _, err := w.Write([]byte("OK"))
		 if err != nil {
			 log.Errorf("Failed to write health check response: %v", err)
		 }
	 })
	 
	 log.Infof("Starting metrics server on %s", serviceConfig.Port)
	 
	 server := &http.Server{
		 Addr:    ":" + serviceConfig.Port,
		 Handler: mux,
	 }
	 
	 if err := server.ListenAndServe(); err != nil {
		 if err != http.ErrServerClosed {
			 log.Errorf("Metrics server error: %v", err)
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