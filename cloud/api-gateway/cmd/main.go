package main

import (
	"github.com/ukama/ukamaX/cloud/api-gateway/cmd/version"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg"
	"github.com/ukama/ukamaX/cloud/api-gateway/pkg/rest"
	"os"

	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
)

var svcConf = pkg.NewConfig()
var ServiceName = "api-gateway"

func main() {
	ccmd.ProcessVersionArgument(ServiceName, os.Args, version.Version)
	initConfig()

	var authMiddleware rest.AuthMiddleware
	authMiddleware = rest.NewKratosAuthMiddleware(&svcConf.Kratos, svcConf.DebugMode)

	if svcConf.BypassAuthMode && svcConf.DebugMode {
		authMiddleware = rest.NewDebugAuthMiddleware()
	}

	r := rest.NewRouter(svcConf.Port, svcConf.DebugMode, authMiddleware, svcConf.Cors, rest.NewClientsSet(&svcConf.Services))
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(ServiceName, svcConf)
}
