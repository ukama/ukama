package main

import (
	"os"

	"github.com/ukama/ukama/systems/data-plan/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := rest.NewClientsSet(&svcConf.Services)

	  metrics.StartMetricsServer(&svcConf.Metrics)

	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf))
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
