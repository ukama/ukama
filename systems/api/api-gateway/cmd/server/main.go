package main

import (
	"os"

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/common/sql"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	resDB := initDb()

	clientSet := client.NewClientsSet(db.NewResourceRepo(resDB),
		client.NewNetworkClient(svcConf.HttpServices.Network))

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

func initDb() sql.Db {
	log.Infof("Initializing Database")

	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db.Resource{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}
