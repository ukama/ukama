package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/mailer/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/mailer/api-gateway/pkg"
	"github.com/ukama/ukama/systems/mailer/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/mailer/api-gateway/pkg/rest"
)

var svcConf = pkg.NewConfig(pkg.SystemName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	logrus.Infof("Starting %s", pkg.ServiceName)
	am := client.NewMailerManager(svcConf.Auth.AuthServerUrl, 3*time.Second, svcConf.Auth.KetoUrl)
	cs := rest.NewClientsSet(am)
	r := rest.NewRouter(cs, rest.NewRouterConfig(svcConf))
	r.Run()
}

func initConfig() {
	config.LoadConfig(pkg.ServiceName, svcConf)
}
