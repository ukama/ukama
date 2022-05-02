package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukama/ukama/services/common/metrics"
	sr "github.com/ukama/ukama/services/common/srvcrouter"
	"github.com/ukama/ukama/testing/factory/internal"
	"github.com/ukama/ukama/testing/factory/internal/server"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/testing/factory/cmd/version"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/config"
)

func main() {
	ccmd.ProcessVersionArgument(internal.ServiceName, os.Args, version.Version)

	/* Log level */
	logrus.SetLevel(logrus.TraceLevel)

	/*Signal handler for SIGINT or SIGTERM to cancel a context in
	order to clean up and shut down gracefully if Ctrl+C is hit. */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* Signal Handling */
	handleSigterm(func() {
		logrus.Infof("Cleaning %s service.", internal.ServiceName)
		/* Call anything required for clean exit */

		cancel()
	})

	/* config parsig */
	initConfig()

	/* Start the HTTP server. */
	startHTTPServer(ctx)
}

/* Start HTTP server for accepting REST  request */
func startHTTPServer(ctx context.Context) {

	logrus.Tracef("Config is %+v", internal.ServiceConfig)

	ext := make(chan error)
	rs := sr.NewServiceRouter(internal.ServiceConfig.ServiceRouter)

	metrics.StartMetricsServer(&internal.ServiceConfig.Metrics)

	r := server.NewRouter(internal.ServiceConfig, rs)
	go r.Run(ext)

	/* Register service */
	if err := rs.RegisterService(internal.ServiceConfig.ApiIf); err != nil {
		logrus.Errorf("Exiting the %s service registration failed.", internal.ServiceName)
		return
	}

	perr := <-ext
	if perr != nil {
		panic(perr)
	}

	logrus.Infof("Exiting service %s", internal.ServiceName)
}

/* initConfig reads in config file, ENV variables, and flags if set. */
func initConfig() {
	internal.ServiceConfig = internal.NewConfig()
	logrus.Tracef("Config is %+v", internal.ServiceConfig)
	config.LoadConfig(internal.ServiceName, internal.ServiceConfig)
	internal.IsDebugMode = internal.ServiceConfig.DebugMode
}

/* Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting. */
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		handleExit()
		logrus.Infof("Exiting %s.", internal.ServiceName)
		os.Exit(1)
	}()

}
