package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukama/ukama/services/common/metrics"

	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg"
	"github.com/ukama/ukama/services/bootstrap/bootstrap/pkg/server"
	sr "github.com/ukama/ukama/services/common/srvcrouter"

	"github.com/ukama/ukama/services/bootstrap/bootstrap/cmd/version"
	ccmd "github.com/ukama/ukama/services/common/cmd"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/config"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/* Log level */
	logrus.SetLevel(logrus.TraceLevel)

	/*Signal handler for SIGINT or SIGTERM to cancel a context in
	order to clean up and shut down gracefully if Ctrl+C is hit. */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* Signal Handling */
	handleSigterm(func() {
		logrus.Infof("Cleaning bootstrap service.")
		/* Call anything required for clean exit */

		cancel()
	})

	/* config parsig */
	initConfig()

	/* Start the HTTP server. */
	startBootstrapServer(ctx)
}

/* Start HTTP server for accepting bootstrap request */
func startBootstrapServer(ctx context.Context) {

	logrus.Tracef("Config is %+v", serviceConfig)

	ext := make(chan error)
	rs := sr.NewServiceRouter(serviceConfig.ServiceRouter)

	metrics.StartMetricsServer(&serviceConfig.Metrics)

	r := server.NewRouter(serviceConfig, rs)

	go r.Run(ext)

	/* Register service */
	if err := rs.RegisterService(serviceConfig.ApiIf); err != nil {
		logrus.Errorf("Exiting the bootstarp service.")
		return
	}

	perr := <-ext
	if perr != nil {
		panic(perr)
	}

	logrus.Infof("Exiting service %s", pkg.ServiceName)
}

/* initConfig reads in config file, ENV variables, and flags if set. */
func initConfig() {
	serviceConfig = pkg.NewConfig()
	logrus.Tracef("Config is %+v", serviceConfig)
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
		logrus.Infof("Exiting bootstrap service.")
		os.Exit(1)
	}()

}
