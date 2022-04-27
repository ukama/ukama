package main

import (
	"os"

	"github.com/ukama/ukamaX/common/metrics"

	"github.com/ukama/ukamaX/templates/rest-service/pkg"

	"github.com/ukama/ukamaX/templates/rest-service/cmd/version"

	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/templates/rest-service/pkg/server"
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
