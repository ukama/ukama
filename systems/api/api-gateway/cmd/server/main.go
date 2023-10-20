package main

import (
	"os"

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"

	log "github.com/sirupsen/logrus"
	prest "github.com/ukama/ukama/systems/api/api-gateway/pkg/rest"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := client.NewClientsSet(
		rest.NewNetworkClient(svcConf.HttpServices.Network),
		rest.NewPackageClient(svcConf.HttpServices.Package),
		rest.NewSubscriberClient(svcConf.HttpServices.Subscriber),
		rest.NewSimClient(svcConf.HttpServices.Sim),
		rest.NewNodeClient(svcConf.HttpServices.Node),
	)

	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}

	router := prest.NewRouter(clientSet, prest.NewRouterConfig(svcConf), ac.AuthenticateUser)
	router.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, svcConf)
}
