package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/hub/distributor/cmd/version"
	"github.com/ukama/ukama/systems/hub/distributor/pkg"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/distribution"
	"github.com/ukama/ukama/systems/hub/distributor/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/*Signal handler for SIGINT or SIGTERM to cancel a context in
	order to clean up and shut down gracefully if Ctrl+C is hit. */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* Signal Handling */
	handleSigterm(func() {
		log.Infof("Cleaning distribution service.")
		/* Call anything required for clean exit */

		cancel()
	})

	/* config parsig */
	initConfig()

	/* Log level */
	log.SetLevel(log.DebugLevel)

	/* Intilaize credentials */
	pkg.InitStoreCredentialsOptions(&serviceConfig.Storage)

	/* Start the HTTP server for chunk distribution */
	go startDistributionServer(ctx)

	/* Start the HTTP server for chunking request. */
	startChunkRequestServer(ctx)
}

/* Start HTTP distribution server for distributing chunks */
func startDistributionServer(ctx context.Context) {
	err := distribution.RunDistribution(ctx, &serviceConfig.Distribution)
	if err != nil {
		log.Errorf("Error while starting distribution server : %s", err.Error())
		os.Exit(1)
	}
}

/* Start HTTP server for accepting chinking request from UkamaHub */
func startChunkRequestServer(ctx context.Context) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		instanceId = uuid.NewV4().String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	r := server.NewRouter(serviceConfig, mbClient)

	metrics.StartMetricsServer(serviceConfig.Metrics)

	go msgBusListener(mbClient)

	r.Run()
}

/* initConfig reads in config file, ENV variables, and flags if set. */
func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

/* Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting. */
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		handleExit()
		log.Infof("Exiting distribution service.")
		os.Exit(1)
	}()
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
