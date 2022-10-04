package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"

	"github.com/ukama/ukama/systems/registry/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/registry/api-gateway/pkg"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	metrics.StartMetricsServer(&svcConf.Metrics)
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
