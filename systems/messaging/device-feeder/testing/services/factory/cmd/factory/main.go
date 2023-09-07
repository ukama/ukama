package main

import (
	"context"
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	sig "github.com/ukama/ukama/systems/common/signal"
	sr "github.com/ukama/ukama/systems/common/srvcrouter"
	"github.com/ukama/ukama/testing/services/factory/internal"
	"github.com/ukama/ukama/testing/services/factory/internal/server"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/testing/services/factory/cmd/version"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
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
	sig.HandleSigterm(func() {
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
