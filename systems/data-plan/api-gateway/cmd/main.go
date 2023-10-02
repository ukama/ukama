package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg"
	"github.com/ukama/ukama/systems/data-plan/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/providers"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := rest.NewClientsSet(&svcConf.Services)
	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		logrus.Errorf("Failed to create auth client: %v", err)
	}
	metrics.StartMetricsServer(&svcConf.Metrics)
	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf), ac.AuthenticateUser)
	r.Run()

}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
