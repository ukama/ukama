package main

import (
	"os"

	"github.com/ukama/ukamaX/common/metrics"

	"github.com/ukama/ukamaX/metrics/node-metrics/pkg"

	"github.com/ukama/ukamaX/metrics/node-metrics/cmd/version"

	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/metrics/node-metrics/pkg/server"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	m, err := pkg.NewMetrics(serviceConfig.NodeMetrics)
	if err != nil {
		panic("Error creating NodeMetrics. Error: " + err.Error())
	}

	r := server.NewRouter(&serviceConfig.Server, m)
	metrics.StartMetricsServer(serviceConfig.Metrics)
	r.Run()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}
