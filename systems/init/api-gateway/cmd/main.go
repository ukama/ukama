package main

import (
	"os"

	"github.com/ukama/ukama/services/common/metrics"

	"github.com/ukama/ukama/systems/init/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg"
	"github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
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
