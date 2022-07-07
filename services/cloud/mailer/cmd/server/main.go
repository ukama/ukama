package main

import (
	"os"

	"github.com/ukama/ukama/services/cloud/mailer/pkg/http"
	sig "github.com/ukama/ukama/services/common/signal"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/mailer/cmd/version"
	"github.com/ukama/ukama/services/cloud/mailer/pkg"
	"github.com/ukama/ukama/services/common/metrics"

	"github.com/num30/config"
	ccmd "github.com/ukama/ukama/services/common/cmd"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	// init config
	initConfig()

	/*Signal handler for SIGINT or SIGTERM to cancel a context in
	order to clean up and shut down gracefully if Ctrl+C is hit. */
	sig.HandleSigterm(func() {
		logrus.Infof("Cleaning %s service.", pkg.ServiceName)
		/* Call anything required for clean exit */
	})

	// start metrics server
	metrics.StartMetricsServer(serviceConfig.Metrics)

	// start init and start queue listener
	m := pkg.NewMail(serviceConfig.Smtp, serviceConfig.TemplatesPath)
	mailer, err := pkg.NewMailer(serviceConfig.Queue, m)
	if err != nil {
		logrus.Fatalf("failed to create mailer: %v", err)
	}

	err = mailer.Start()
	if err != nil {
		logrus.Fatalf("failed to start mailer: %v", err)
	}

	// start http server
	r := http.NewRouter(serviceConfig.Server)
	r.Run()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = &pkg.Config{}
	reader := config.NewConfReader(pkg.ServiceName)
	err := reader.Read(serviceConfig)
	if err != nil {
		logrus.Fatalf("Failed to read config: %v", err)
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}
