package main

import (
	"os"

	"github.com/ukama/ukama/services/cloud/api-gateway/cmd/version"
	"github.com/ukama/ukama/services/cloud/api-gateway/pkg"
	"github.com/ukama/ukama/services/cloud/api-gateway/pkg/client"
	"github.com/ukama/ukama/services/cloud/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := rest.NewClientsSet(&svcConf.Services)

	var authMiddleware rest.AuthMiddleware
	authMiddleware = rest.NewKratosAuthMiddleware(&svcConf.Kratos,
		client.NewRegistry(svcConf.Services.Network, svcConf.Services.Org, svcConf.Services.Node, svcConf.Services.Timeout), svcConf.DebugMode)

	if svcConf.BypassAuthMode && svcConf.DebugMode {
		authMiddleware = rest.NewDebugAuthMiddleware()
	}

	r := rest.NewRouter(authMiddleware, clientSet, rest.NewRouterConfig(svcConf))
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
