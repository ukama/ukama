package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukama/ukama/services/common/metrics"

	"github.com/ukama/ukamaX/hub/distributor/pkg"
	"github.com/ukama/ukamaX/hub/distributor/pkg/distribution"
	"github.com/ukama/ukamaX/hub/distributor/pkg/server"

	"github.com/ukama/ukamaX/hub/distributor/cmd/version"

	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
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
		logrus.Infof("Cleaning distribution service.")
		/* Call anything required for clean exit */

		cancel()
	})

	/* config parsig */
	initConfig()

	/* Log level */
	logrus.SetLevel(logrus.DebugLevel)

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
		logrus.Errorf("Error while starting distribution server : %s", err.Error())
		os.Exit(1)
	}
}

/* Start HTTP server for accepting chinking request from UkamaHub */
func startChunkRequestServer(ctx context.Context) {
	r := server.NewRouter(serviceConfig)
	metrics.StartMetricsServer(&serviceConfig.Metrics)
	r.Run()
}

/* initConfig reads in config file, ENV variables, and flags if set. */
func initConfig() {
	serviceConfig = pkg.NewConfig()
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
		logrus.Infof("Exiting distribution service.")
		os.Exit(1)
	}()

}
