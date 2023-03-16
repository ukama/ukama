package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"

	"github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg"
	"github.com/ukama/ukama/systems/subscriber/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	subscriberConfig()

	clientSet := rest.NewClientsSet(&svcConf.Services)

	metrics.StartMetricsServer(&svcConf.Metrics)

	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf))
	r.Run()
}

func subscriberConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
