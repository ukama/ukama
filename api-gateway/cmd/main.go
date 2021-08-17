package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/config"
	"ukamaX/api-gateway/pkg"
	"ukamaX/api-gateway/pkg/rest"
)

var svcConf = pkg.NewConfig()
var ServiceName = "api-gateway"

func main() {
	log.Infof("Starting " + ServiceName)
	initConfig()

	var authMiddleware rest.AuthMiddleware
	authMiddleware = rest.NewKratosAuthMiddleware(&svcConf.Kratos, svcConf.DebugMode)

	if svcConf.BypassAuthMode && svcConf.DebugMode {
		authMiddleware = rest.NewDebugAuthMiddleware()
	}

	r := rest.NewRouter(svcConf.Port, svcConf.DebugMode, authMiddleware, rest.NewClientsSet(&svcConf.Services))
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(ServiceName, svcConf)
}
