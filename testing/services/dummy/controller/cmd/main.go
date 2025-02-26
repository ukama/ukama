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
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"

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
    controllerServer := server.NewControllerServer(serviceConfig.OrgName)

    grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
        generated.RegisterMetricsControllerServer(s, controllerServer)
    })

    go grpcServer.StartServer()

    // Start the Prometheus metrics server
   // Start the Prometheus metrics server
go func() {
    // Create a new HTTP server mux
    mux := http.NewServeMux()
    
    // Register the Prometheus handler
    mux.Handle("/metrics", promhttp.Handler())
    
    // Add a basic health check endpoint
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    // Log that we're starting the metrics server
    address := ":2112"
    log.Infof("Starting metrics server on %s", address)
    
    // Start the server and handle errors properly
    server := &http.Server{
        Addr:    address,
        Handler: mux,
    }
    
    if err := server.ListenAndServe(); err != nil {
        if err != http.ErrServerClosed {
            log.Errorf("Metrics server error: %v", err)
        }
    }
}()
    // Set up graceful shutdown - use only one signal handler
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    
    // Wait for termination signal
    sig := <-sigs
    log.Infof("Received signal %v, shutting down...", sig)
    
    // Clean up resources
    controllerServer.Cleanup()
    
    log.Infof("Exiting service %s", pkg.ServiceName)
}
 
 func waitForExit() {
	 sigs := make(chan os.Signal, 1)
	 signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	 done := make(chan bool, 1)
	 
	 go func() {
		 sig := <-sigs
		 log.Info(sig)
		 done <- true
	 }()
 
	 log.Debug("awaiting terminate/interrupt signal")
	 <-done
	 log.Infof("exiting service %s", pkg.ServiceName)
 }