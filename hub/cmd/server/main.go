package main

import (
	"github.com/ukama/ukamaX/common/metrics"
	"os"
	"time"

	"github.com/ukama/ukamaX/hub/pkg"

	"github.com/ukama/ukamaX/hub/cmd/version"

	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/hub/pkg/server"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	storage := pkg.NewMinioWrapper(&serviceConfig.Storage)

	r := server.NewRouter(&serviceConfig.Server, storage, time.Duration(serviceConfig.Storage.TimeoutSecond)*time.Second)
	metrics.StartMetricsServer(serviceConfig.Metrics)
	r.Run()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}
