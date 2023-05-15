package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
)

var svcConf = pkg.NewConfig(pkg.SystemName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	logrus.Infof("Starting %s", pkg.ServiceName)

	orgRegistryClient, err := client.NewOrgMemberRoleClient(svcConf.Auth.OrgUrl, pkg.IsDebugMode)
	if err != nil {
		logrus.Fatalf("orgRegistry Client initilization failed. Error: %v", err.Error())
	}
	am := client.NewAuthManager(svcConf.Auth.AuthServerUrl, 3*time.Second,orgRegistryClient)
	cs := rest.NewClientsSet(am)
	r := rest.NewRouter(cs, rest.NewRouterConfig(svcConf, svcConf.AuthKey))
	r.Run()
}

func initConfig() {
	config.LoadConfig(pkg.ServiceName, svcConf)
}
