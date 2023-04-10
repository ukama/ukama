package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/metrics"

	"github.com/ukama/ukama/systems/metrics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	log.Infof("Config %+v", svcConf.MetricsConfig)
	clientSet := rest.NewClientsSet(&svcConf.Services, svcConf.MetricsStore, svcConf.DebugMode)

	metrics.StartMetricsServer(&svcConf.MetricsServer)

	m, err := pkg.NewMetrics(svcConf.MetricsConfig)
	if err != nil {
		panic("Error creating NodeMetrics. Error: " + err.Error())
	}

	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf), m)
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
