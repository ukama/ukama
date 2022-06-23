package main

import (
	"context"
	"os"

	"github.com/ukama/ukama/services/cloud/notify/internal"
	"github.com/ukama/ukama/services/cloud/notify/internal/db"
	"github.com/ukama/ukama/services/cloud/notify/internal/server"
	"github.com/ukama/ukama/services/common/metrics"
	sig "github.com/ukama/ukama/services/common/signal"
	"github.com/ukama/ukama/services/common/sql"

	"github.com/ukama/ukama/services/cloud/notify/cmd/version"
	ccmd "github.com/ukama/ukama/services/common/cmd"

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
	sig.HandleSigterm(func() {
		logrus.Infof("Cleaning %s service.", internal.ServiceName)
		/* Call anything required for clean exit */

		cancel()
	})

	/* Config parsing */
	initConfig()

	/*  Database */
	d := initDb()

	/* Start the HTTP server. */
	startHTTPServer(ctx, d)
}

func initDb() sql.Db {
	logrus.Infof("Initializing Database")
	d := sql.NewDb(internal.ServiceConfig.DB, internal.ServiceConfig.DebugMode)
	err := d.Init(&db.Notification{})
	if err != nil {
		logrus.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

/* Start HTTP server for accepting REST  request */
func startHTTPServer(ctx context.Context, d sql.Db) {

	logrus.Tracef("Config is %+v", internal.ServiceConfig)

	metrics.StartMetricsServer(&internal.ServiceConfig.Metrics)

	r := server.NewRouter(internal.ServiceConfig, db.NewNotificationRepo(d))

	r.Run()

	logrus.Infof("Exiting service %s", internal.ServiceName)
}

/* initConfig reads in config file, ENV variables, and flags if set. */
func initConfig() {
	internal.ServiceConfig = internal.NewConfig()
	logrus.Tracef("Config is %+v", internal.ServiceConfig)
	config.LoadConfig(internal.ServiceName, internal.ServiceConfig)
	internal.IsDebugMode = internal.ServiceConfig.DebugMode
}
