package main

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

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
	am := client.NewAuthManager(svcConf.Auth.AuthServerUrl, 3*time.Second)
	cs := rest.NewClientsSet(am)
	r := rest.NewRouter(cs, rest.NewRouterConfig(svcConf, svcConf.AuthKey))
	r.Run()
}

func initConfig() {
	config.LoadConfig(pkg.ServiceName, svcConf)
}
