package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/notification/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg"
	"github.com/ukama/ukama/systems/notification/api-gateway/pkg/rest"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initializeNotificationConfig()

	clientSet := rest.NewClientsSet(&svcConf.Services)
	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}

	router := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf), ac.AuthenticateUser)
	router.Run()
}

func initializeNotificationConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
