package main

import (
	"os"
	"time"

	"github.com/ukama/ukama/services/common/metrics"

	"github.com/ukama/ukama/services/hub/hub/pkg"

	"github.com/ukama/ukama/services/hub/hub/cmd/version"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/hub/hub/pkg/server"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	storage := pkg.NewMinioWrapper(&serviceConfig.Storage)
	chunker := pkg.NewChunker(&serviceConfig.Chunker, storage)

	r := server.NewRouter(&serviceConfig.Server, storage, chunker,
		time.Duration(serviceConfig.Storage.TimeoutSecond)*time.Second)
	metrics.StartMetricsServer(serviceConfig.Metrics)
	r.Run()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}
