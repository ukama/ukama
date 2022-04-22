package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ukama/openIoR/services/common/metrics"
	sr "github.com/ukama/openIoR/services/common/srvcrouter"
	"github.com/ukama/openIoR/services/factory/nmr/internal/db"
	"github.com/ukama/openIoR/services/factory/nmr/pkg"
	"github.com/ukama/openIoR/services/factory/nmr/pkg/server"

	ccmd "github.com/ukama/openIoR/services/common/cmd"
	"github.com/ukama/openIoR/services/common/sql"
	"github.com/ukama/openIoR/services/factory/nmr/cmd/version"

	"github.com/sirupsen/logrus"
	"github.com/ukama/openIoR/services/common/config"
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

	d := initDb()

	/* Start the HTTP server. */
	startBootstrapServer(ctx, d)
}

func initDb() sql.Db {
	logrus.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Node{}, &db.Module{})
	if err != nil {
		logrus.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

/* Start HTTP server for accepting bootstrap request */
func startBootstrapServer(ctx context.Context, d sql.Db) {

	logrus.Tracef("Config is %+v", serviceConfig)

	rs := sr.NewServiceRouter(serviceConfig.ServiceRouter)

	/* Register service */
	if err := rs.RegisterService(serviceConfig.ApiIf); err != nil {
		logrus.Errorf("Exiting the bootstarp service.")
		//return
	}

	metrics.StartMetricsServer(&serviceConfig.Metrics)

	r := server.NewRouter(serviceConfig, rs, db.NewNodeRepo(d), db.NewModuleRepo(d))
	r.Run()
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
