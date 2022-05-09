package main

import (
	"os"

	"github.com/ukama/ukama/services/common/metrics"

	"github.com/ukama/ukama/services/metrics/node-metrics/pkg"

	"github.com/ukama/ukama/services/metrics/node-metrics/cmd/version"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/metrics/node-metrics/pkg/server"
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
	// validate
	for k := range serviceConfig.NodeMetrics.RawQueries {
		if _, ok := serviceConfig.NodeMetrics.Metrics[k]; ok {
			panic("Duplicate metric name: " + k)
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}
