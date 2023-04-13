package main

import (
	"net/http"
	"net/http/cookiejar"
	"os"

	ory "github.com/ory/client-go"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/auth/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg"
	"github.com/ukama/ukama/systems/auth/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	crest "github.com/ukama/ukama/systems/common/rest"
)

var svcConf = pkg.NewConfig(pkg.SystemName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	rc, err := crest.NewRestClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		logrus.Errorf("Can't conncet to %v url. Error %v", svcConf.Auth.AuthServerUrl, err.Error())
	}
	svcConf.R = rc
	svcConf.R.C = svcConf.R.C.SetBaseURL(svcConf.Auth.AuthServerUrl)
	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{
			URL: svcConf.Auth.AuthServerUrl,
		},
	}
	jar, _ := cookiejar.New(nil)
	oc := ory.NewAPIClient(configuration)
	oc.GetConfig().HTTPClient = &http.Client{
		Jar: jar,
	}
	r := rest.NewRouter(rest.NewRouterConfig(svcConf, oc, svcConf.AuthKey))
	r.Run()
}

func initConfig() {
	config.LoadConfig(pkg.ServiceName, svcConf)
}
