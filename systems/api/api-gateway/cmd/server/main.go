package main

import (
	"os"

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := client.NewClientsSet(
		client.NewNetworkClient(svcConf.HttpServices.Network),
		client.NewPackageClient(svcConf.HttpServices.Package),
		client.NewSubscriberClient(svcConf.HttpServices.Subscriber),
		client.NewSimClient(svcConf.HttpServices.Sim),
		client.NewNodeClient(svcConf.HttpServices.Node),
	)

	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}

	router := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf), ac.AuthenticateUser)
	router.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, svcConf)
}
