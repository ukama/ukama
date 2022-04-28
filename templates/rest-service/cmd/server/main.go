package main

import (
	"os"

	"github.com/ukama/ukama/services/common/metrics"

	"github.com/ukama/ukama/services/templates/rest-service/pkg"

	"github.com/ukama/ukama/services/templates/rest-service/cmd/version"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/templates/rest-service/pkg/server"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	r := server.NewRouter(&serviceConfig.Server)
	metrics.StartMetricsServer(serviceConfig.Metrics)
	r.Run()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}
